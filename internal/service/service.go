package service

import (
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	log  *zap.Logger
	repo *repository.UrlRepository
	r    Repository
}

var (
	siteNotLoad        = errors.New("site data not loaded. Please wait")
	siteNotConnect     = errors.New("site is unavailable")
	siteNameNotCorrect = errors.New("incorrect site requested")
)

type Repository interface {
	ReadSiteInfo(s string) model.SiteResponseInfo
	ReadMinResponseSite() model.SiteResponseInfo
	ReadMaxResponseSite() model.SiteResponseInfo
	ReadCountSiteRequest(s string) (uint64, error)
	ReadCountMaxRequest() uint64
	ReadCountMinRequest() uint64
}

func NewService(log *zap.Logger, repo *repository.UrlRepository) *Service {
	log.Debug("Register service...")
	return &Service{
		log:  log,
		repo: repo}
}

func (s Service) GetSiteInfo(site string) (model.SiteResponseInfo, error) {
	s.log.Debug("Get site response...")
	result := s.repo.ReadSiteInfo(site)
	if result.Code == -1 {
		return result, siteNameNotCorrect
	}
	s.repo.WriteCountSiteRequest(site)

	if result.Code == 0 {
		return result, siteNotLoad
	}
	if result.Code != http.StatusOK {
		return result, siteNotConnect
	}
	return result, nil
}

func (s Service) GetSiteMinResponse() (model.SiteResponseInfo, error) {
	s.log.Debug("Get Min site response...")
	result := s.repo.ReadMinResponseSite()
	s.repo.WriteCountMinRequest()
	if result.Code == 0 {
		return result, siteNotLoad
	}
	if result.Code != http.StatusOK {
		return result, siteNotConnect
	}
	return result, nil
}

func (s Service) GetSiteMaxResponse() (model.SiteResponseInfo, error) {
	s.log.Debug("Get Max site response...")
	result := s.repo.ReadMaxResponseSite()
	s.repo.WriteCountMaxRequest()
	if result.Code == 0 {
		return result, siteNotLoad
	}
	if result.Code != http.StatusOK {
		return result, siteNotConnect
	}
	return result, nil
}

func (s Service) GetSiteStat(site string) (uint64, error) {
	s.log.Debug(fmt.Sprintf("Get site %s stat...", site))
	result, err := s.repo.ReadCountSiteRequest(site)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s Service) GetMinStat() uint64 {
	s.log.Debug("Get Min site stat...")
	return s.repo.ReadCountMinRequest()
}

func (s Service) GetMaxStat() uint64 {
	s.log.Debug("Get Max site stat...")
	return s.repo.ReadCountMaxRequest()
}
