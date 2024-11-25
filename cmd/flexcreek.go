package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)

	fmt.Println("Running Flex Creek on port 8080")

	http.ListenAndServe(":8080", mux)
}

// index handler
func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Flex Creek"))
}
