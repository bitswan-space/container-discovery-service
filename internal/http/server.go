package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bitswan.space/container-discovery-service/internal/logger"
	"bitswan.space/container-discovery-service/internal/models"
	"bitswan.space/container-discovery-service/internal/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

const (
	ShutdownTimeout = 1 * time.Second
	ReadTimeout     = 5 * time.Second
	WriteTimeout    = 10 * time.Second
	IdleTimeout     = 120 * time.Second
	DefaultAddr     = ":8080"
)

type HttpServer struct {
	router *chi.Mux

	CDSRepository repository.CDSRepository
}

// NewHttpServer initializes a new HTTP server with a router.
func NewHttpServer() *HttpServer {
	router := chi.NewRouter()

	return &HttpServer{
		router: router,
	}
}

// Router sets up the router with all routes, separating protected from unprotected.
func (s *HttpServer) Router() {
	// Apply CORS middleware to all routes.
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		Response(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	s.router.Get("/api/dashboard-entries", s.FetchDashboardEntries)
	s.router.Post("/api/dashboard-entries", s.CreateDashboardEntry)
	s.router.Delete("/api/dashboard-entries/{id}", s.DeleteDashboardEntry)
	s.router.Put("/api/dashboard-entries/{id}", s.UpdateDashboardEntry)

}

// Run starts the HTTP server on a specified port.
func (s *HttpServer) Run() {
	port := os.Getenv("PORT")
	logger.Info.Printf("Running HTTP server on port %s", port)
	logger.Info.Printf("%v", http.ListenAndServe(fmt.Sprintf(":%s", port), s.router))
}

// FetchDashboardEntries fetches the dashboard entries from the json config file.
func (s *HttpServer) FetchDashboardEntries(w http.ResponseWriter, r *http.Request) {

	topology, err := s.CDSRepository.FetchDashboardEntries(r.Context())
	if err != nil {
		Error(w, r, err)
		return
	}

	// Respond with the fetched topology.
	Response(w, http.StatusOK, topology)
}

func (s *HttpServer) CreateDashboardEntry(w http.ResponseWriter, r *http.Request) {

	entry := &models.DashboardEntry{}
	if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
		Error(w, r, err)
		return
	}

	dashboardEntry, err := s.CDSRepository.CreateDashboardEntry(r.Context(), entry)
	if err != nil {
		Error(w, r, err)
		return
	}

	Response(w, http.StatusCreated, dashboardEntry)
}

func (s *HttpServer) DeleteDashboardEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.CDSRepository.DeleteDashboardEntry(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	Response(w, http.StatusNoContent, nil)
}

func (s *HttpServer) UpdateDashboardEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entry := &models.DashboardEntry{}
	if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
		Error(w, r, err)
		return
	}

	dashboardEntry, err := s.CDSRepository.UpdateDashboardEntry(r.Context(), id, entry)
	if err != nil {
		Error(w, r, err)
		return
	}

	Response(w, http.StatusOK, dashboardEntry)
}
