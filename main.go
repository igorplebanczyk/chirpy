package main

import (
	"go-server/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
	db             *database.Database
}

func main() {
	db, err := database.NewDatabase("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		db:             db,
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    "192.168.1.27:8080",
		Handler: mux,
	}

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetMetricsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpHandlerSpecific)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	err = server.ListenAndServe()
	if err != nil {
		return
	}
}
