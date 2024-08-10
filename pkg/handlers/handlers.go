package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"os"

	"go-whyye/pkg/db"
	"go-whyye/pkg/quote"
	"go-whyye/pkg/user"
)

type Response struct {
	Message string `json:"message"`
}

type UsersListAllBody struct {
	Users []user.User `json:"users"`
}

func UsersListAllHandler(w http.ResponseWriter, r *http.Request) {
	database, err := db.NewDatabase()
	if err != nil {
		http.Error(w, "Error creating database connection", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	users, err := database.UsersListAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&UsersListAllBody{
		Users: users,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type QuotesListAllBody struct {
	Quotes []quote.Quote `json:"quotes"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type UserByIdQuotesBody struct {
	Quotes []quote.Quote `json:"quotes"`
	UserID int           `json:"userId"`
	Limit  int           `json:"limir"`
	Offset int           `json:"offset"`
}

func QuotesListAllHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	var userId int
	if qs.Get("user_id") != "" {
		id, err := strconv.Atoi(qs.Get("user_id"))
		if err != nil || id <= 0 {
			log.Println(err)
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
		userId = id
	}

	var limit int
	if qs.Get("limit") != "" {
		qLimit, err := strconv.Atoi(qs.Get("limit"))
		if err != nil || qLimit <= 0 {
			log.Println(err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		limit = qLimit
	} else {
		limit = 100
	}

	var offset int
	if qs.Get("offset") != "" {
		qOffset, err := strconv.Atoi(qs.Get("offset"))
		if err != nil || qOffset < 0 {
			log.Println(err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset = qOffset
	} else {
		offset = 0
	}

	database, err := db.NewDatabase()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating database connection", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var quotes []quote.Quote
	var quoteErr error
	if userId > 0 && limit > 1 {
		quotes, quoteErr = database.UserByIdQuotes(userId, limit, offset)
	} else {
		quotes, quoteErr = database.QuotesListAll(limit, offset)
	}

	if quoteErr != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve quotes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if userId > 0 {
		err = json.NewEncoder(w).Encode(&UserByIdQuotesBody{
			Quotes: quotes,
			UserID: userId,
			Limit:  limit,
		})
	} else {
		err = json.NewEncoder(w).Encode(&QuotesListAllBody{
			Quotes: quotes,
			Limit:  limit,
		})
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	qs := r.URL.Query()
	var userId int
	if qs.Get("user_id") != "" {
		id, err := strconv.Atoi(qs.Get("user_id"))
		if err != nil || id <= 0 {
			log.Println(err)
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		userId = id
	}

	database, err := db.NewDatabase()
	if err != nil {
		http.Error(w, "Error creating database connection", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	users, err := database.UsersListAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	var qt quote.Quote
	if userId > 0 {
		quotes, err := database.UserByIdQuotes(userId, 100, 0)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to get user quote", http.StatusInternalServerError)
			return
		}

		qt = quote.GetRandom(quotes)
	}
	var AnalyticsSrc string
	var AnalyticsId string
	var AnalyticsHost string

	app_env, exists := os.LookupEnv("APP_ENV")
	if exists && app_env == "production" {
		AnalyticsSrc = "https://umami.littlevibe.net/script.js"
		AnalyticsId = "aa24ac89-ef05-4bdf-a99b-402988a6f226"
		AnalyticsHost = "https://umami.littlevibe.net/script.js"
	}

	vars := map[string]interface{}{
		"Users":  users,
		"UserId": userId,
		"Quote":  qt.String(),
		"Year":   time.Now().Year(),

		"AnalyticsSrc": AnalyticsSrc,
		"AnalyticsId": AnalyticsId,
		"AnalyticsHost": AnalyticsHost,
	}

	err = tmpl.ExecuteTemplate(w, "index.html", vars)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

