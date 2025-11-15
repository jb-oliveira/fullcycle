package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// bookHotelMain()
	serverContext()
}

func serverContext() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request Iniciada")
	defer log.Println("Request Finalizada")
	select {
	case <-time.After(5 * time.Second):
		log.Println("Request processada com sucesso")
		w.Write([]byte("Request processada com sucesso"))
	case <-ctx.Done():
		log.Println("Request cancelada")
		http.Error(w, "Request cancelada", http.StatusInternalServerError)
	}
}

func bookHotelMain() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {

	select {
	case <-ctx.Done():
		fmt.Println("Cancelado. Timeout reached")
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Hotel booked.")
	}
}
