package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) requestCount(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetCount(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getProfaneWords() [3]string { return [...]string{"kerfuffle", "sharbert", "fornax"} }

func cleanChrip(chirp string) string {
	cleanedWords := []string{}
	for word := range strings.SplitSeq(chirp, " ") {
		isProfaneWord := false
		for _, profaneWord := range getProfaneWords() {
			if profaneWord == strings.ToLower(word) {
				isProfaneWord = true
			}
		}
		if isProfaneWord {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " ")
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	if err := decoder.Decode(&req); err != nil {
		respondWithErrorJSON(w, 500, err)
		return
	}
	if len(req.Body) > 140 {
		respondWithErrorJSON(w, 400, fmt.Errorf("Chirp is too long"))
		return
	}
	res := responseBody{CleanedBody: cleanChrip(req.Body)}
	respondWithJSON(w, 200, res)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Run() {
	cfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.requestCount)
	mux.HandleFunc("POST /admin/reset", cfg.resetCount)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
