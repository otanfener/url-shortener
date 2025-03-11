package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/otanfener/url-shortener/internal/server/transport"
	"github.com/otanfener/url-shortener/internal/service/dto"
)

type Handler struct {
	service service
}

func NewHandler(s service) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Routes(router chi.Router) {
	router.Post("/urls", h.ShortenURL)
	router.Get("/{code}", h.RedirectURL)
}
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req transport.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	dtoReq := dto.ShortenRequest{LongURL: req.LongURL}
	shortCode, err := h.service.ShortenURL(dtoReq)
	if err != nil {
		http.Error(w, "failed to shorten url", http.StatusInternalServerError)
		return
	}

	resp := transport.ShortenURLResponse{ShortCode: shortCode}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	longURL, err := h.service.RedirectURL(shortCode)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
