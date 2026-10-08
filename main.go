package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	database "github.com/Jonak-Adipta-Kalita/JAK-RSS-Aggregator-Project/src/database/generated"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT isn't found in Environment")
	}
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("DB_URL isn't found in Environment")
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
		return
	}

	apiCfg := ApiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJson(w, 200, struct{}{})
	})
	v1Router.Post("/user", apiCfg.handleCreateUser)
	v1Router.Get("/user", apiCfg.handleGetUserByAPIKey)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	log.Printf("Server starting on port %v", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
