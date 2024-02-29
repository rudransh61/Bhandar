package main

import (
	"sync"
)

// Database represents a simple in-memory key-value store.
type Database struct {
	data     map[string]string
	mutex    sync.Mutex
	notifier chan string
}

// NewDatabase creates a new instance of the database.
func NewDatabase() *Database {
	return &Database{
		data:     make(map[string]string),
		notifier: make(chan string),
	}
}

// Set adds or updates a key-value pair in the database and notifies subscribers.
func (db *Database) Set(key, value string) {
	db.mutex.Lock()
	db.data[key] = value
	db.mutex.Unlock()
	db.notifier <- key
}

// Get retrieves the value associated with a key from the database.
func (db *Database) Get(key string) (string, bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	value, exists := db.data[key]
	return value, exists
}

// Subscribe adds a new channel to receive real-time updates.
func (db *Database) Subscribe(ch chan string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.notifier = ch
}

// Unsubscribe removes a channel from receiving real-time updates.
func (db *Database) Unsubscribe() {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	close(db.notifier)
	db.notifier = make(chan string)
}
