package main

import (
	"log"
	"log_analyse/DB"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	envPkg "github.com/joho/godotenv"
)

func main() {
	envPkg.Load(".env")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbProp := DB.SDB{}

	dbInitErr := dbProp.Init()

	if dbInitErr != nil {
		os.Exit(1)
	}

	router := chi.NewRouter()

	var srv *http.Server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		var health bool = dbProp.Health()

		if !health {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database is not healthy"))
			return
		}
		w.Write([]byte("Database is healthy"))
	})

	log.Printf("Server started on port %s", port)

	var err = srv.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server: ", err)
		dbProp.Close()
	}

}
