package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var accessCount uint64 = 0

func main() {

	// with mutex
	// m := sync.Mutex{}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	m.Lock()
	// 	accessCount++
	// 	m.Unlock()
	// 	w.Write([]byte(fmt.Sprintf("Access Count: %d", accessCount)))
	// })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&accessCount, 1)
		w.Write([]byte(fmt.Sprintf("Access Count: %d", accessCount)))
	})

	http.ListenAndServe(":3000", nil)
}
