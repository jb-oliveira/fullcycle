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

func BuscaDolarHandler(w http.ResponseWriter, r *http.Request) {
	response, err := loadDataFromUrl()
	if err != nil {
		log.Println("Error lendo url", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dolar, err := strconv.ParseFloat(response.Bid, 64)
	if err != nil {
		log.Println("Error convertendo string to float64:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("Erro criando arquivo", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = f.WriteString(fmt.Sprintf("Dólar: %f", dolar))
	if err != nil {
		log.Println("Erro escrevendo arquivo", err)
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
	http.HandleFunc("/", BuscaDolarHandler)
	http.ListenAndServe(":8081", nil)
}
