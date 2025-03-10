package handler

import (
	"encoding/json"
	"net/http"

	"github.com/otanfener/url-shortener/internal/handler/transport"
	"github.com/otanfener/url-shortener/internal/service/dto"
)

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req transport.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	dtoReq := dto.ShortenRequest{LongURL: req.LongURL}
	shortCode, err := h.shortenURLService.ShortenURL(dtoReq)
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
