package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DOLAR_API_URL  = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	API_TIMEOUT_MS = 200
	DB_TIMEOUT_MS  = 10
)

type (
	USDBRL struct {
		gorm.Model
		ID         int64  `json:"-" gorm:"primaryKey,autoIncrement"`
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	}
	ApiResponse struct {
		USDBRL USDBRL `json:"USDBRL"`
	}
)

func GetDollarPrice(ctx context.Context) (*ApiResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DOLAR_API_URL, nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("request timed out, insuficient time to make request")
			return nil, fmt.Errorf("request timed out")
		}

		return nil, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		strBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(strBody))
	}

	var apiResponse ApiResponse
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func InitServer(ctx context.Context) *http.Server {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&USDBRL{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), API_TIMEOUT_MS*time.Millisecond)
		defer cancel()

		price, err := GetDollarPrice(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel = context.WithTimeout(r.Context(), DB_TIMEOUT_MS*time.Millisecond)
		defer cancel()

		if err := db.Create(&price.USDBRL).Error; err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("database operation timed out, insuficient time to save data")
				http.Error(w, "database operation timed out", http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(price.USDBRL.Bid))
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// listen and server in a goroutine
	// when the context is canceled, the server will stop
	go func() {
		log.Println("Server running on port 8080")

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server exited: %v", err)
		}
	}()

	return &srv
}
