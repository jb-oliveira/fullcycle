package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	appDB *gorm.DB
)

func GetDB() *gorm.DB {
	return appDB
}

func InitializeSQLite() (*gorm.DB, error) {
	fmt.Println("Initializing SQLite database")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&USDBRL{})
	if err != nil {
		return nil, err
	}
	fmt.Println("SQLite database initialized")

	return db, nil
}

type USDBRL struct {
	gorm.Model
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

type UrlRequest struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type UrlResponse struct {
	Bid string `json:"bid"`
}

func BuscaDolarHandler(w http.ResponseWriter, r *http.Request) {
	request, err := loadDataFromUrl()
	if err != nil {
		log.Println("Error consulting url", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = saveToDataBase(request)
	if err != nil {
		log.Println("Error saving to DB", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := UrlResponse{
		Bid: request.USDBRL.Bid,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func saveToDataBase(entity *UrlRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	return gorm.G[USDBRL](GetDB()).Create(ctx, &entity.USDBRL)
}

func loadDataFromUrl() (*UrlRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	var result UrlRequest
	err = json.Unmarshal(body, &result)
	return &result, nil
}

func main() {
	var err error
	appDB, err = InitializeSQLite()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/cotacao", BuscaDolarHandler)
	http.ListenAndServe(":8080", nil)
}
