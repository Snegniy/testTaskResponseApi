package service

import (
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.Logger
	cfg  *config.Config
	repo *repository.UrlRepository
}

func NewService(log *zap.Logger, cfg *config.Config) *Service {
	log.Debug("Register service...")
	r := repository.NewRepository(log, cfg)
	go response(r, cfg.UrlRepo.Timeout, cfg.UrlRepo.Refresh)
	return &Service{log, cfg, r}
}

type Services interface {
	ReadSite(site string) (float64, error)
	//AddCount(r Repository) (*UrlRepository, error)
	GetMinMaxSite(m string) (float64, string, error)
}

func (s *Service) GetSite(site string) (float64, error) {
	s.log.Debug("Get Min Site response...")
	time, err := s.repo.ReadSite(site)
	if err != nil {
		return 0, err
	}
	return time, nil
}

func (s *Service) GetMin() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	time, site, err := s.repo.GetMinMaxSite("min")
	if err != nil {
		return 0, "", err
	}
	return time, site, nil
}

func (s *Service) GetMax() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	return 0, "", nil
}

func (s *Service) GetStat() (float64, string, error) {
	s.log.Debug("Get Min Site response...")
	return 0, "", nil
}
