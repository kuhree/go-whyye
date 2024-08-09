package main

import (
	"fmt"
	"net/http"
	"os"

	"go-whyye/pkg/db"
	"go-whyye/pkg/handlers"
)

func main() {
	err := db.PrepareDb()
	if err != nil {
		panic(err)
	}

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/api/users", handlers.UsersListAllHandler)
	http.HandleFunc("/api/quotes", handlers.QuotesListAllHandler) // support ?userId=[id] to filter by user
	http.HandleFunc("/", handlers.IndexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not set
	}

	fmt.Printf("Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
