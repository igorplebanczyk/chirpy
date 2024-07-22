package main

import (
	"net/http"
)

func main() {
	apiCfg := apiConfig{
		fileServerHits: 0,
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
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
