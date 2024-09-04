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

type Request struct {
	Numbers []int `json:"numbers"`
}

type Response struct {
	Result int `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid input. Unable to parse JSON."})
		return
	}

	log.Printf("Parsed request: %+v", req)

	if len(req.Numbers) == 0 {
		log.Println("Empty numbers array")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid input. Please provide an array of integers in the \"numbers\" field."})
		return
	}

	sum := 0
	for _, num := range req.Numbers {
		sum += num
	}

	log.Printf("Calculated sum: %d", sum)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: sum})
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
