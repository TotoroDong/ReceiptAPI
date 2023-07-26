package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        string  `json:"total"`
	Items        []Items `json:"items"`
}

type Items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var receiptsMap = make(map[string]int)

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var receipt Receipt
	json.Unmarshal(body, &receipt)

	id := uuid.New().String()

	points := calculatePoints(receipt)
	receiptsMap[id] = points

	resp := IDResponse{ID: id}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

	fmt.Printf("Processed receipt with ID: %s\n", id) //Check for ID in map
}

type PointsIDResponse struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/") // split the path
	id := pathSegments[2]                          // the ID should be the third segment
	fmt.Printf("Requested points for ID: %s\n", id)
	points, exists := receiptsMap[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	resp := PointsIDResponse{ID: id, Points: points}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1
	points += len(regexp.MustCompile(`\w`).FindAllString(receipt.Retailer, -1))

	// Rule 2
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalFloat == math.Round(totalFloat) {
		points += 50
	}

	// Rule 3
	if math.Remainder(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4
	points += (len(receipt.Items) / 2) * 5

	// Rule 5
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Round(price * 0.2))
		}
	}

	// Rule 6
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() <= 16 {
		points += 10
	}

	return points
}

func main() {
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
