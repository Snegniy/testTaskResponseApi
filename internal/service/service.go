package service

import (
	"errors"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.Logger
	repo *repository.UrlRepository
}

var (
	siteNotLoad    = errors.New("site data not loaded. Please wait")
	siteNotConnect = errors.New("site is unavailable")
)

type Services interface {
	ReadSiteInfo(s string) (model.SiteResponseInfo, error)
	ReadMinResponseSite() (model.SiteResponseInfo, error)
	ReadMaxResponseSite() (model.SiteResponseInfo, error)
	ReadCountSiteRequest(s string) (uint64, error)
	ReadCountMaxRequest() uint64
	ReadCountMinRequest() uint64
}

func NewService(log *zap.Logger, repo *repository.UrlRepository) *Service {
	log.Debug("Register service...")
	return &Service{log, repo}
}

func (s Service) GetSiteInfo(site string) (model.SiteResponseInfo, error) {
	s.log.Debug("Get Site response...")
	result, err := s.repo.ReadSiteInfo(site)
	if err != nil {
		return result, err
	}
	if result.ResponseTime == 0.0 {
		if result.Code == 0 {
			return result, siteNotLoad
		}
		if result.Code != 200 {
			return result, siteNotConnect
		}
	}
	return result, nil
}

func (s Service) GetSiteMinResponse() (model.SiteResponseInfo, error) {
	s.log.Debug("Get Min Site response...")
	result, _ := s.repo.ReadMinResponseSite()
	if result.Code == 0 {
		return result, siteNotLoad
	}
	if result.Code != 200 {
		return result, siteNotConnect
	}

	return result, nil
}

func (s Service) GetSiteMaxResponse() (model.SiteResponseInfo, error) {
	s.log.Debug("Get Max Site response...")
	result, _ := s.repo.ReadMaxResponseSite()
	if result.Code == 0 {
		return result, siteNotLoad
	}
	if result.Code != 200 {
		return result, siteNotConnect
	}
	return result, nil
}

func (s Service) GetSiteStat(site string) (uint64, error) {
	s.log.Debug("Get Site stat...")
	result, err := s.repo.ReadCountSiteRequest(site)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s Service) GetMinStat() uint64 {
	s.log.Debug("Get Min Site stat...")
	return s.repo.ReadCountMinRequest()
}

func (s Service) GetMaxStat() uint64 {
	s.log.Debug("Get Max Site stat...")
	return s.repo.ReadCountMaxRequest()
}
