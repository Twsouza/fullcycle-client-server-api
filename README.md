# Client Server Api - Fullcycle Course Challenge

This project is implemented in Go and consists of two systems: a client (client.go) and a server (server.go). These systems interact to retrieve and store currency exchange rates between USD and BRL (Brazilian Real).

Client (client.go): Sends an HTTP request to the server to fetch the current USD to BRL exchange rate and saves the received rate in a text file.
Server (server.go): Consumes an external API to get the exchange rate, logs the data into a SQLite database, and returns the exchange rate to the client.

This project demonstrates the use of HTTP web servers, contexts, databases, and file manipulation in Go.

## How to run

To run the project, you need to have Go installed on your machine.

Run the server:

```bash
go run main.go
```
