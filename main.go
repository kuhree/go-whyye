package main

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"

	"go-whyye/pkg/services"
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

func main() {
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
