package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func LogError(err error) {
	log.SetPrefix("ERROR: ")                             // Add a prefix to log entries
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Include date, time, and file info
	log.Printf("An error occurred: %v", err)
}

func main() {

	executeGet()
	executePost()
	customRequest()
	contextExample()

}

func contextExample() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://www.google.com", nil)
	if err != nil {
		LogError(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		LogError(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		LogError(err)
		return
	}
	println(string(body))
}

func customRequest() {
	c := http.Client{}
	req, err := http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		LogError(err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		LogError(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogError(err)
		return
	}
	println(string(body))

}

func executeGet() {
	c := http.Client{Timeout: time.Duration(1) * time.Microsecond}
	res, err := c.Get("http://www.google.com")
	if err != nil {
		LogError(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		LogError(err)
		return
	}
	println(string(body))
}

func executePost() {
	c := http.Client{Timeout: time.Duration(1) * time.Microsecond}

	jsonVar := bytes.NewBuffer([]byte(`{ "name":"Joao"}`))

	res, err := c.Post("http://localhost:8080/post", "application/json", jsonVar)

	if err != nil {
		LogError(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		LogError(err)
		return
	}
	println(string(body))
}
