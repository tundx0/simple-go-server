package main

import (
	"encoding/json"
	"fmt"
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

func sumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid input. Unable to parse JSON."})
		return
	}

	if len(req.Numbers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid input. Please provide an array of integers in the \"numbers\" field."})
		return
	}

	sum := 0
	for _, num := range req.Numbers {
		sum += num
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Result: sum})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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