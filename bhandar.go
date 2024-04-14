package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// CacheEntry represents a key-value pair with an expiration time.
type CacheEntry struct {
	value      string
	expiration time.Time
}

// Cache represents an in-memory key-value cache with TTL support.
type Cache struct {
	data map[string]CacheEntry // Internal data store
	mu   sync.RWMutex          // Mutex for thread safety
}

// NewCache creates a new instance of Cache.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheEntry),
	}
}

// Set adds or updates a key-value pair in the cache with a TTL.
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(ttl)
	c.data[key] = CacheEntry{value: value, expiration: expiration}
}

// Get retrieves the value associated with a key from the cache.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok {
		return "", false
	}
	if time.Now().After(entry.expiration) {
		// Entry has expired
		delete(c.data, key)
		return "", false
	}
	return entry.value, true
}

// TTL returns the remaining time to live for a given key.
func (c *Cache) TTL(key string) (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok {
		return 0, false
	}
	ttl := entry.expiration.Sub(time.Now())
	if ttl <= 0 {
		// Entry has expired
		delete(c.data, key)
		return 0, false
	}
	return ttl, true
}

// List returns a list of all keys in the cache.
func (c *Cache) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func main() {
	// Create a new cache instance
	cache := NewCache()

	// Start HTTP server
	go func() {
		http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
			key := r.FormValue("key")
			value := r.FormValue("value")
			ttl, err := time.ParseDuration(r.FormValue("ttl")) // Parse TTL from request parameters
			if err != nil {
				http.Error(w, "Invalid TTL", http.StatusBadRequest)
				return
			}

			cache.Set(key, value, ttl)
			fmt.Fprintf(w, "Key-value pair (%s, %s) set successfully with TTL %s\n", key, value, ttl)
		})

		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			key := r.FormValue("key")

			value, found := cache.Get(key)
			if found {
				fmt.Fprintf(w, "%s", value)
			} else {
				fmt.Fprintf(w, "(nil)")
			}
		})

		http.HandleFunc("/ttl", func(w http.ResponseWriter, r *http.Request) {
			key := r.FormValue("key")

			ttl, found := cache.TTL(key)
			if found {
				fmt.Fprintf(w, "TTL for key %s: %s\n", key, ttl)
			} else {
				fmt.Fprintf(w, "Key not found: %s\n", key)
			}
		})

		http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
			keys := cache.List()
			fmt.Fprintf(w, "Keys in cache: %s\n", strings.Join(keys, ", "))
		})

		http.ListenAndServe(":8080", nil)
	}()

	// Real-time command line interface
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("127.0.0.1:8000>")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			parts := strings.Fields(input)
			if len(parts) == 0 {
				continue
			}

			command := parts[0]
			switch command {
			case "set":
				if len(parts) != 4 {
					fmt.Println("Invalid set command. Format: set key value ttl")
					continue
				}
				key := parts[1]
				value := parts[2]
				ttlStr := parts[3]
				ttl, err := time.ParseDuration(ttlStr)
				if err != nil {
					fmt.Println("Invalid TTL duration:", err)
					continue
				}
				set(key, value, ttl)
			case "get":
				if len(parts) != 2 {
					fmt.Println("Invalid get command. Format: get key")
					continue
				}
				key := parts[1]
				get(key)
			case "ttl":
				if len(parts) != 2 {
					fmt.Println("Invalid ttl command. Format: ttl key")
					continue
				}
				key := parts[1]
				ttl(key)
			case "list":
				list()
			case "exit":
				fmt.Println("Exiting...")
				os.Exit(0)
			default:
				fmt.Println("Invalid command")
			}
		}
	}()

	// Block main goroutine to keep the application running
	select {}
}

func set(key, value string, ttl time.Duration) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/set?key=%s&value=%s&ttl=%s", key, value, ttl.String()))
	if err != nil {
		fmt.Println("Error setting key-value pair:", err)
		return
	}
	defer resp.Body.Close()

	// fmt.Println("Response:", resp.Status)
}

func get(key string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/get?key=%s", key))
	if err != nil {
		fmt.Println("Error getting value by key:", err)
		return
	}
	defer resp.Body.Close()

	// fmt.Println("Response:")
	// fmt.Println(resp.Status)
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
}

func ttl(key string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/ttl?key=%s", key))
	if err != nil {
		fmt.Println("Error getting TTL for key:", err)
		return
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
}

func list() {
	resp, err := http.Get("http://localhost:8080/list")
	if err != nil {
		fmt.Println("Error getting list of keys:", err)
		return
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
}
