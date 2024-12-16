package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/imhasandl/go-restapi/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	filepath := "."
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Set Port in env")
	}
	dbURl := os.Getenv("DB_URL")
	if dbURl == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURl)
	if err != nil {
		log.Fatal("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepath)))

	mux.HandleFunc("/status", apiCfg.handlerStatusCheck)

	mux.HandleFunc("users", apiCfg.handlerUserCreate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server running on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
