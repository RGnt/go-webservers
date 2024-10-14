package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsPrint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//"<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>\n"
	status := fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %d times!</p>\n</body>\n</html>\n", cfg.fileserverHits)
	w.Write([]byte(status))
}

func (cfg *apiConfig) middlewareMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Chirp struct {
		Body string `json:"body"`
	}
	type RetErr struct {
		Error string `json:"error"`
	}
	type RetValid struct {
		Valid bool `json:"valid"`
	}

	retChirpIsValid := RetValid{
		Valid: true,
	}

	errorSomethingWentWrong := RetErr{
		Error: "Something went wrong",
	}

	errorChirpTooLong := RetErr{
		Error: "Chirp is too long",
	}
	decoder := json.NewDecoder(r.Body)
	body := Chirp{}
	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(500)
		dat, _ := json.Marshal(errorSomethingWentWrong)
		w.Write(dat)
		return
	}
	if len(body.Body) > 140 {
		w.WriteHeader(400)
		dat, err := json.Marshal(errorChirpTooLong)
		if err != nil {
			return
		}
		w.Write(dat)
		return
	}

	chirpValid, err := json.Marshal(retChirpIsValid)
	if err != nil {
		w.WriteHeader(400)
		dat, _ := json.Marshal(errorSomethingWentWrong)
		w.Write(dat)
		return
	}
	w.WriteHeader(200)
	w.Write(chirpValid)
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{fileserverHits: 0}
	handler := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsPrint)
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
