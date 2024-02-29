package main

import (
	"fmt"
	"net/http"
	// "example.go"
)

func main() {
	db := NewDatabase()

	// Create a subscriber channel
	updates := make(chan string, 10)
	db.Subscribe(updates)

	// Start a goroutine to listen for updates
	go func() {
		for key := range updates {
			fmt.Printf("Real-Time Update: Key '%s' changed\n", key)
		}
	}()

	// API endpoints
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value, exists := db.Get(key)
		if exists {
			fmt.Fprintf(w, "Key: %s, Value: %s\n", key, value)
		} else {
			http.Error(w, "Key not found", http.StatusNotFound)
		}
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		db.Set(key, value)
		fmt.Fprintf(w, "Key '%s' set to value '%s'\n", key, value)
	})

	// Run the web server
	fmt.Println("Web API running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
