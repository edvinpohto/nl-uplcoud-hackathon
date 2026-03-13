package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/data", dataHandler)
	mux.HandleFunc("/api/slow", slowHandler)
	mux.HandleFunc("/", rootHandler)

	log.Printf("victim-app listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"hostname":  hostname(),
	})
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	items := make([]map[string]interface{}, 10)
	for i := range items {
		items[i] = map[string]interface{}{
			"id":    i + 1,
			"value": rand.Intn(1000),
			"ts":    time.Now().UTC().UnixMilli(),
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":    items,
		"hostname": hostname(),
	})
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate a slow endpoint (50-200ms)
	delay := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(delay)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"delay_ms":   delay.Milliseconds(),
		"hostname":   hostname(),
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"service": "victim-app",
		"version": "1.0.0",
	})
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
