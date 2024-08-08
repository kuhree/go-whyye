package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go-whyye/pkg/services"
	"go-whyye/pkg/db"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleKanye(w http.ResponseWriter, r *http.Request) {
	baseUrl := "https://api.kanye.rest"
	svc := services.NewKanyeRestSvc(baseUrl)

	quote, err := svc.FetchQuote()
	if err != nil {
		fmt.Println(err)
		return
	}

	data := Response{
		Status:  "ok",
		Message: quote.Quote,
	}

	json.NewEncoder(w).Encode(data)
}

func prepareDb() (*db.Database, error) {
	db, err := db.NewDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.CreateSchema()
	if err != nil {
			return nil, err
	}

	err = db.SeedDatabase()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	_, err := prepareDb()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/kanye", handleKanye)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not set
	}

	fmt.Printf("Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
