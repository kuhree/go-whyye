package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-whyye/lib"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleKanye(w http.ResponseWriter, r *http.Request) {
	baseUrl := "https://api.kanye.rest"
	svc := lib.NewKanyeRestSvc(baseUrl)

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

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
