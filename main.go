package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func logRequest(r *http.Request) {
	log.Println("--- New Request ---")
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL)
	log.Printf("Headers: %v", r.Header)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	// Restore the body for later use
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Printf("Body: %s", string(body))
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var numbers []int
	err := json.NewDecoder(r.Body).Decode(&numbers)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid input. Unable to parse JSON.", http.StatusBadRequest)
		return
	}

	log.Printf("Parsed numbers: %v", numbers)

	if len(numbers) == 0 {
		log.Println("Empty numbers array")
		http.Error(w, "Invalid input. Please provide an array of integers.", http.StatusBadRequest)
		return
	}

	sum := 0
	for _, num := range numbers {
		sum += num
	}

	log.Printf("Calculated sum: %d", sum)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, sum)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprint(w, "Service is running")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/", healthCheckHandler)

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
