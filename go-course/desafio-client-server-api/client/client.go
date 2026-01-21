package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UrlResponse struct {
	Bid string `json:"bid"`
}

func SearchDollarHandler(w http.ResponseWriter, r *http.Request) {
	response, err := loadDataFromUrl()
	if err != nil {
		log.Println("Error reading url", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dollar, err := strconv.ParseFloat(response.Bid, 64)
	if err != nil {
		log.Println("Error converting string to float64:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("Error creating file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.WriteString(fmt.Sprintf("Dollar: %f", dollar))
	if err != nil {
		log.Println("Error writing file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func loadDataFromUrl() (*UrlResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result UrlResponse
	err = json.Unmarshal(body, &result)
	return &result, nil
}

func main() {
	http.HandleFunc("/", SearchDollarHandler)
	http.ListenAndServe(":8081", nil)
}
