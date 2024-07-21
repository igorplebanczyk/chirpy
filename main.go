package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    "192.168.1.27:8080",
		Handler: mux,
	}
	mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", readinessHandler)

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
