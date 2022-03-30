package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	JsonContent     = "application/json"
	NotFoundMessage = "page not found"
	OK              = "ok"
	DefaultEndpoint = "/"
	HealthEndpoint  = "/health"
)

type DefaultResponse struct {
	App string `json:"app"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type HealthResponse struct {
	Status string `json:"status"`
}

type Server struct {
	http.Handler
	app string
}

func NewServer(app string, router *http.ServeMux, HealthHandler http.HandlerFunc) *Server {
	s := new(Server)
	s.Handler = router
	s.app = app

	if nil == HealthHandler {
		HealthHandler = StubHealthHandler
	}

	router.HandleFunc(DefaultEndpoint, s.DefaultHandler)
	router.HandleFunc(HealthEndpoint, HealthHandler)

	return s
}

func (s *Server) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, JsonContent)

	if "/" == r.URL.Path {
		HandleResponse(w, DefaultResponse{s.app}, http.StatusOK)
	} else {
		HandleError(w, NewNotFoundError(errors.New(NotFoundMessage)))
	}
}

// StubHealthHandler default always healthy handler.
// Must be implemented per app with proper checks.
func StubHealthHandler(w http.ResponseWriter, r *http.Request) {
	HandleResponse(w, HealthResponse{OK}, http.StatusOK)
}

func HandleResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set(ContentType, JsonContent)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func HandleError(w http.ResponseWriter, err CodeError) {
	w.Header().Set(ContentType, JsonContent)
	w.WriteHeader(err.Code())
	json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
}
