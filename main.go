package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"go-whyye/pkg/db"
	"go-whyye/pkg/handlers"
)

func main() {
	err := db.PrepareDb()
	if err != nil {
		panic(err)
	}

	app_env, exists := os.LookupEnv("APP_ENV")
	if !exists {
		app_env = "development"
	}

	if sentry_dsn, exists := os.LookupEnv("SENTRY_DSN"); app_env == "production" && exists {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentry_dsn,
			TracesSampleRate: 1.0,
		})

		if err != nil {
			log.Fatalf("Failed sentry.Init: %s", err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{ Repanic: true, })
	http.HandleFunc("/sentry", sentryHandler.HandleFunc(func(rw http.ResponseWriter, r *http.Request) {
		panic("y tho")
	}))

	http.Handle("/static/", sentryHandler.Handle(http.StripPrefix("/static", http.FileServer(http.Dir("./static")))))
	http.HandleFunc("/api/users", sentryHandler.HandleFunc(handlers.UsersListAllHandler))
	http.HandleFunc("/api/quotes", sentryHandler.HandleFunc(handlers.QuotesListAllHandler)) // support ?userId=[id] to filter by user
	http.HandleFunc("/", sentryHandler.HandleFunc(handlers.IndexHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not set
	}

	fmt.Printf("Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
