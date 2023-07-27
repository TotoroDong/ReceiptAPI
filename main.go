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

// For debugging purpose. Point breakdown.
type ReceiptPoints struct {
	RetailerNamePoints      int
	RoundDollarPoints       int
	QuarterMultiplePoints   int
	ItemCountPoints         int
	ItemDescriptionPoints   int
	OddDayPoints            int
	AfternoonPurchasePoints int
	TotalPoints             int
}

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

// In memory storage
var receiptsMap = make(map[string]ReceiptPoints)

// Function that handle incoming receipts.
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var receipt Receipt
	json.Unmarshal(body, &receipt)

	//Generate Unqiue ID for each receipt using uuid
	id := uuid.New().String()

	points := calculatePoints(receipt)
	receiptsMap[id] = points

	resp := IDResponse{ID: id}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

	// Debugging purpose: Print out points for each rule
	fmt.Println("Receipt ID:", id)
	fmt.Println("Points from retailer name:", points.RetailerNamePoints)
	fmt.Println("Points from round dollar total:", points.RoundDollarPoints)
	fmt.Println("Points from quarter multiple total:", points.QuarterMultiplePoints)
	fmt.Println("Points from item count:", points.ItemCountPoints)
	fmt.Println("Points from item description:", points.ItemDescriptionPoints)
	fmt.Println("Points from odd day:", points.OddDayPoints)
	fmt.Println("Points from afternoon purchase:", points.AfternoonPurchasePoints)
	fmt.Println("Total points:", points.TotalPoints)
}

type PointsIDResponse struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

// Process the request to obtain Point base on id
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/") // split the path
	id := pathSegments[2]                          // the ID should be the third segment
	fmt.Printf("Requested points for ID: %s\n", id)
	receiptPoints, exists := receiptsMap[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	resp := PointsIDResponse{ID: id, Points: receiptPoints.TotalPoints}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// Calculate Points for each receipt and store it in-memory
func calculatePoints(receipt Receipt) ReceiptPoints {
	points := ReceiptPoints{}

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points.RetailerNamePoints = len(regexp.MustCompile(`\w`).FindAllString(receipt.Retailer, -1))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalFloat == math.Round(totalFloat) {
		points.RoundDollarPoints = 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Remainder(totalFloat, 0.25) == 0 {
		points.QuarterMultiplePoints = 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points.ItemCountPoints = (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		if len(item.ShortDescription)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			itemPoints := int(math.Ceil(price * 0.2))
			points.ItemDescriptionPoints += itemPoints

			// Print the item description and points
			fmt.Printf("Item Description: %s, Points Earned: %d\n", item.ShortDescription, itemPoints)
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate) //2006-01-02 Time package format
	if purchaseDate.Day()%2 != 0 {
		points.OddDayPoints = 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime) //15:04 time package format
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() <= 16 {
		points.AfternoonPurchasePoints = 10
	}

	// Calculate total points
	points.TotalPoints = points.RetailerNamePoints + points.RoundDollarPoints + points.QuarterMultiplePoints + points.ItemCountPoints + points.ItemDescriptionPoints + points.OddDayPoints + points.AfternoonPurchasePoints

	return points
}

func main() {
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
