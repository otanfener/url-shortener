package service

import (
	"errors"
	"log"

	"github.com/otanfener/url-shortener/internal/service/dto"
	"github.com/otanfener/url-shortener/pkg/base62"
)

type storage interface {
	SaveURLMapping(shortCode, longURL string) error
	GetLongURL(shortCode string) (string, error)
	CheckShortCodeExists(shortCode string) (bool, error)
}
type counter interface {
	NextID() (int64, error)
}

type Service struct {
	storage storage
	counter counter
}

func NewService(s storage, c counter) *Service {
	return &Service{storage: s, counter: c}

}

func (s *Service) ShortenURL(req dto.ShortenRequest) (string, error) {

	if req.LongURL == "" {
		return "", errors.New("long URL is required")
	}

	var shortCode string

	id, err := s.counter.NextID()
	if err != nil {
		log.Printf("failed to get next ID from counter: %v", err)
		return "", err
	}
	shortCode = base62.Encode(id)

	// Save URL mapping
	if err := s.storage.SaveURLMapping(shortCode, req.LongURL); err != nil {
		log.Printf("failed to save URL mapping in storage: %v", err)
		return "", err
	}
	return shortCode, nil
}

func (s *Service) RedirectURL(code string) (string, error) {
	return s.storage.GetLongURL(code)
}
