package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CLIENT_TIMEOUT_MS = 300
)

func GetDollarPrice(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, CLIENT_TIMEOUT_MS*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("request timed out, insuficient time to make request")
			return
		}

		log.Fatalf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		strBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("unexpected status code %v: %v", res.StatusCode, err)
		}

		log.Fatal(string(strBody))
	}

	log.Println("request made successfully")

	// save the response body to cotacao.txt
	strBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	response := fmt.Sprintf("Dolar: %v", string(strBody))
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(response)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}

	log.Println("response saved to cotacao.txt")
}
