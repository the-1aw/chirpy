package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/the-1aw/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
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

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithErrorJSON(w, http.StatusForbidden, fmt.Errorf("Cannot reset database if not in dev mode"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(r.Context())
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
		respondWithErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	if len(req.Body) > 140 {
		respondWithErrorJSON(w, http.StatusBadRequest, fmt.Errorf("Chirp is too long"))
		return
	}
	res := responseBody{CleanedBody: cleanChrip(req.Body)}
	respondWithJSON(w, http.StatusOK, res)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Run() error {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		db: dbQueries,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("POST /api/users", getCreateUserHandler(&cfg))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.requestCount)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	return server.ListenAndServe()
}
