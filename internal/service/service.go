package service

import (
	"errors"
	"fmt"

	"github.com/otanfener/url-shortener/internal/domain"
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
type logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}
type Service struct {
	storage storage
	counter counter
	logger  logger
}

func NewService(s storage, c counter, l logger) *Service {
	return &Service{storage: s, counter: c, logger: l}
}

func (s *Service) ShortenURL(req dto.ShortenRequest) (string, error) {
	if req.LongURL == "" {
		return "", fmt.Errorf("%w: long URL is required", domain.ErrInvalidInput)
	}
	var shortCode string

	id, err := s.counter.NextID()
	if err != nil {
		s.logger.Error("failed to generate next ID", map[string]interface{}{"error": err.Error()})
		return "", fmt.Errorf("%w: failed to generate next ID", domain.ErrCounterFailure)
	}
	shortCode = base62.Encode(id)

	// Save URL mapping
	err = s.storage.SaveURLMapping(shortCode, req.LongURL)
	if err != nil {
		if errors.Is(err, domain.ErrStorageFailure) {
			s.logger.Error("failed to save url mapping", map[string]interface{}{
				"error":      err.Error(),
				"short_code": shortCode,
				"long_url":   req.LongURL},
			)
			return "", err
		}
		return "", fmt.Errorf("%w: failed to save URL mapping", domain.ErrInternal)
	}
	return shortCode, nil
}
func (s *Service) RedirectURL(code string) (string, error) {
	return s.storage.GetLongURL(code)
}
