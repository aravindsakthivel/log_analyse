package main

import (
	"log"
	"log_analyse/DB"
	"log_analyse/api"
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

	routeHandler := api.SRouter{}

	routeHandler.RoutesInit(router)

	log.Printf("Server started on port %s", port)

	var err = srv.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server: ", err)
		dbProp.Close()
	}

}
