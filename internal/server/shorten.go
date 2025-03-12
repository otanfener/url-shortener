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
	encoder *Encoder
}

func NewHandler(s service, e *Encoder) *Handler {
	return &Handler{service: s, encoder: e}
}
func (h *Handler) Routes(router chi.Router) {
	router.Post("/urls", h.ShortenURL)
	router.Get("/{code}", h.RedirectURL)
}
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req transport.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.encoder.ErrorResponse(ctx, w, err, http.StatusBadRequest)
		return
	}

	dtoReq := dto.ShortenRequest{LongURL: req.LongURL}
	shortCode, err := h.service.ShortenURL(dtoReq)
	if err != nil {
		h.encoder.ErrorResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	resp := transport.ShortenURLResponse{ShortCode: shortCode}
	h.encoder.StatusResponse(ctx, w, resp, http.StatusCreated)
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	ctx := r.Context()
	longURL, err := h.service.RedirectURL(shortCode)
	if err != nil {
		h.encoder.ErrorResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
	h.encoder.RedirectResponse(ctx, w, longURL, http.StatusFound)
}
