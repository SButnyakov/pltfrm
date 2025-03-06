package service

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/models"
	"url-shortener/pkg/utils"
)

type urlRepository interface {
	Create(url *models.URL) error
	GetByURL(url string) (*models.URL, error)
	GetByAddress(address string) (*models.URL, error)
}

type urlService struct {
	urlRepository urlRepository
}

func NewUrlService(urlRepository urlRepository) *urlService {
	return &urlService{urlRepository: urlRepository}
}

func (s *urlService) Create(address string) (string, error) {
	m, err := s.urlRepository.GetByAddress(address)
	if err == nil {
		return m.URL, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	url, err := s.generateUniqueURL()
	if err != nil {
		return "", err
	}

	model := &models.URL{
		Address: address,
		URL:     url,
	}

	return url, s.urlRepository.Create(model)
}

func (s *urlService) GetByURL(url string) (*models.URL, error) {
	return s.urlRepository.GetByURL(url)
}

func (s *urlService) generateUniqueURL() (string, error) {
	const maxAttempts = 100
	for i := 0; i < maxAttempts; i++ {
		url := utils.GenerateRandomString(5)
		_, err := s.urlRepository.GetByURL(url)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return url, nil
			}
			return "", err
		}
	}
	return "", fmt.Errorf("failed to generate unique URL after %d attempts", maxAttempts)
}
