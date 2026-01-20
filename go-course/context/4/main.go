package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("Server stopped.")

}

func longRunningTask(ctx context.Context, r *http.Request) (string, error) {
	resultChan := make(chan string, 1)

	go func() {
		fmt.Println("Executando")
		timeStr := r.URL.Query().Get("time")
		timeInt, err := strconv.Atoi(timeStr)
		if err != nil {
			resultChan <- err.Error()
		}
		time.Sleep(time.Duration(timeInt) * time.Second) // Simulate a long-running operation
		resultChan <- "Task completed"
	}()

	select {
	case <-ctx.Done(): // Context was canceled or timed out
		return "", ctx.Err()
	case result := <-resultChan: // Task completed successfully
		return result, nil
	}
}

func teste(r *http.Request) error {
	timeStr := r.URL.Query().Get("time")
	timeInt, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Tempo em segundos: %d", timeInt)
	time.Sleep(time.Duration(timeInt) * time.Second)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// newFunction(ctx, r, w)
	result, err := longRunningTask(ctx, r)
	if err != nil {
		log.Printf("error: %v", err)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("result: %v", result)
	w.Write([]byte(result))
}

func newFunction(ctx context.Context, r *http.Request, w http.ResponseWriter) {
	log.Println("Request Iniciada")
	defer log.Println("Request Finalizada")
	select {
	case <-time.After(3 * time.Second):
		log.Println("Entrou")
		err := teste(r)
		if err != nil {
			log.Println("Request processada com erro")
			w.Write([]byte("Request processada com erro"))
			return
		}
		log.Println("Request processada com sucesso")
		w.Write([]byte("Request processada com sucesso"))
	case <-ctx.Done():
		log.Println("Request cancelada")
		http.Error(w, "Request cancelada", http.StatusInternalServerError)
	}
}
