package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/imhasandl/go-restapi/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db        *database.Queries
	status    string
	jwtSecret string
}

func main() {
	filepath := "."
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Set Port in env")
	}
	dbURl := os.Getenv("DB_URL")
	if dbURl == "" {
		log.Fatal("DB_URL must be set")
	}
	status := os.Getenv("STATUS")
	if status == "" {
		log.Fatal("Please set your status")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable should be set")
	}

	dbConn, err := sql.Open("postgres", dbURl)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:        dbQueries,
		status:    status,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepath)))

	mux.HandleFunc("GET /status", apiCfg.handlerStatusCheck)

	mux.HandleFunc("POST /api/users/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/users/register", apiCfg.handlerUserCreate)

	mux.HandleFunc("GET /api/users", apiCfg.handlerListAllUsers)
	mux.HandleFunc("GET /api/users/", apiCfg.handlerGetUserByEmail)
	mux.HandleFunc("GET /api/users/{user_id}", apiCfg.handlerGetUserByID)

	mux.HandleFunc("POST /api/posts", apiCfg.handlerCreatePost)
	mux.HandleFunc("GET /api/posts", apiCfg.handlerListPosts)

	mux.HandleFunc("GET /api/posts/{post_id}", apiCfg.hanlerGetPostByID)
	mux.HandleFunc("PUT /api/posts/{post_id}", apiCfg.handlerChangePostByID)
	mux.HandleFunc("DELETE /api/posts/{post_id}", apiCfg.handlerDeletePostByID)

	mux.HandleFunc("DELETE /admin/reset/users", apiCfg.handlerResetUsers)
	mux.HandleFunc("DELETE /admin/reset/posts", apiCfg.handlerResetPosts)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}

	fmt.Printf("Server running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
