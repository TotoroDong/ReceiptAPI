# Fetch Receipt Reward Challenge

The project is created for Fetch receipt reward challenege. [Github](https://github.com/fetch-rewards/receipt-processor-challenge)

This is a mini-project developed using Go programming language. Utilized various libraries to facilitate the development.

# Pre-requisites
You should have Go installed on your system. If not, please follow the instructions provided [here](https://go.dev/doc/install) to install Go.

For Visual Studio Code please install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)


You would also need to install the package [gopls](https://pkg.go.dev/golang.org/x/tools/gopls#section-readme) and [dlv](https://github.com/go-delve/delve). Visual Studio Code should help yout install, but
if not then you can also run the comment in terminal `go install -v goland.org/x/tools/gopls@latest` and `go install github.com/go-delve/delve/cmd/dlv@latest`

If you run into any other issues check the [Docs](https://github.com/golang/vscode-go/wiki/tools)

The code Utilized external libraries such as google uuid. If you don't have mod file in the folder run `go mod init mymodule` to create one and then use `Go get github.com/google/uuid` to add that library/package.


# Running the program
1. Open a terminal and enter the primary folder where "main.go" is located.
2. In the main terminal type `go run main.go`. The program should start running and waiting for json file to be send over.
3. Open a second terminal and enter the sub folder named "JsonSender"
4. In the second terminal type `go run send.go` The program will send a specific json file located within the same folder to the main program.
4a. The default json to send is "git-receipt.json" If you would like to switch, go to line 12 of "send.go" and change the file name.
```go
jsonFilePath := "git-receipt.json" // file path of json file to send
```
Replace "git-receipt.json" to the new json file you want to send. 
Note: This can be done mutitple time since the main program is always running.

# Result
- The moment you use `go run send.go` In the terminal, it should return a random unique ID 
Example:
`{"id":"3a89e386-2bb6-4b23-b727-1b5813587fb4"}`

- After the `id` has returned, check your main terminal that was running `main.go` The terminal or console should display the Receipt ID and the point gain as well as the point structure. These display are solely for debugging purposes and can be removed
- To access the points earned for any receipt, open up a browser and enter the following link. 
**Remember to replace the id with what you got previously
`http://localhost:8080/receipts/id(to be replaced)/points`

Using the id from the above example, the link would become 

`http://localhost:8080/receipts/3a89e386-2bb6-4b23-b727-1b5813587fb4/points`

and the site should return something similar to `{"id":"3a89e386-2bb6-4b23-b727-1b5813587fb4","points":28}`

Note: The main terminal should display `Requested points for ID: 3a89e386-2bb6-4b23-b727-1b5813587fb4` This is also for debugging to see what is going on.

# Future Optimizations
The json file could possibly be send using curl such as 

`curl -X POST -H "Content-Type: application/json" -d @~/Desktop/goproject/receipt.json http://localhost:8080/receipts/process
`

or

`Invoke-WebRequest -Uri "http://localhost:8080/receipts/process" -Method POST -ContentType "application/json" -Body (Get-Content -Path "C:\Users\dc872\Desktop\Goproject\morning-recepit.json" -Raw)
`

The learning process continues!

# Contact
Email: dc872722944@Outlook.com
