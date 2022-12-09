package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	m "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/middleware"
	s "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/storage"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	middlewareStruct = &m.MiddlewareStruct{}
)

type StorageHandlers struct {
	storage s.StorageInterface
}

func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var reader io.Reader

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
		defer r.Body.Close()
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return body, nil
}

func SetValues(filePath string, baseURL string, server string) {
	middlewareStruct.InitMiddlewareStruct(filePath, baseURL, server)
}

func (sh StorageHandlers) PostAddURLHandler(w http.ResponseWriter, r *http.Request) {
	urlBytes, err := ReadBody(w, r)
	if err != nil {
		log.Printf("failed read request: %v", err)
		http.Error(w, "failed read request", http.StatusInternalServerError)
		return
	}
	url := string(urlBytes)
	fullShortenURL, err := sh.storage.AddURL(url, middlewareStruct.FilePath, middlewareStruct.BaseURL)

	if err != nil {
		log.Printf("failed save url '%v': %v", url, err)
		http.Error(w, "failed save url", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fullShortenURL))
}

func (sh StorageHandlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var (
		newURLFull    m.URLFull
		newURLShorten m.URLShorten
	)

	urlBytes, err := ReadBody(w, r)
	if err != nil {
		log.Printf("failed read request: %v", err)
		http.Error(w, "failed read request", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(urlBytes, &newURLFull)
	if err != nil {
		log.Println(err)
		return
	}
	fullShortenURL, err := sh.storage.AddURL(newURLFull.URLFull, middlewareStruct.FilePath, middlewareStruct.BaseURL)

	if err != nil {
		log.Printf("failed save url '%v': %v", string(urlBytes), err)
		http.Error(w, "failed save url", http.StatusInternalServerError)
		return
	}

	newURLShorten.URLShorten = fullShortenURL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newURLShorten)
}

func (sh StorageHandlers) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID parameter must be Integer type", http.StatusBadRequest)
		return
	}

	url := sh.storage.SearchURL(id)

	if url == "" {
		http.Error(w, "There is no url with this id", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(url))
	}

}

func NewRouter(storage s.StorageInterface) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	handlers := StorageHandlers{storage}

	return router.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostAddURLHandler)
		r.Get("/{id}", handlers.GetURLHandler)
		r.Post("/api/shorten", handlers.ShortenHandler)
	})
}
