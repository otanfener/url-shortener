package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	longURL, err := h.shortenURLService.RedirectURL(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
