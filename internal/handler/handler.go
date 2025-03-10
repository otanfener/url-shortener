package handler

import "github.com/otanfener/url-shortener/internal/service/dto"

type shortenURLService interface {
	ShortenURL(req dto.ShortenRequest) (string, error)
	RedirectURL(code string) (string, error)
}

type Handler struct {
	shortenURLService shortenURLService
}

func NewHandler(s shortenURLService) *Handler {
	return &Handler{s}
}
