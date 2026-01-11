package main

import (
	"fmt"
	"log"
	"net/http"

	"tatria/langchain"
	"tatria/route"
)

func main() {
	// Create a new ServeMux
	router := http.NewServeMux()

	route.Routes(router)

	// initialize langchain agent
	err := langchain.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server and listen on port 8080
	fmt.Printf("Server starting on port 8080\n")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
