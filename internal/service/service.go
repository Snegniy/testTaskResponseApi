package service

import (
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.Logger
	repo *repository.UrlRepository
}

type Services interface {
	ReadSite(site string) (float64, error)
	//AddCount(r Repository) (*UrlRepository, error)
	GetMinMaxSite(m string) (float64, string, error)
}

func NewService(log *zap.Logger, repo *repository.UrlRepository) *Service {
	log.Debug("Register service...")
	return &Service{log, repo}
}

func (s Service) GetSite(site string) (float64, error) {
	s.log.Debug("Get Min Site response...")
	time, err := s.repo.ReadSite(site)
	if err != nil {
		return 0, err
	}
	return time, nil
}

func (s Service) GetMin() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	time, site, err := s.repo.GetMinMaxSite("min")
	if err != nil {
		return 0, "", err
	}
	return time, site, nil
}

func (s Service) GetMax() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	return 0, "", nil
}

func (s Service) GetStat() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	return 0, "", nil
}
