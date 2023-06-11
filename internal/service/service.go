package service

import (
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
}

type OutputUserInfo struct {
	RequestName  string `json:"request_name"`
	ResponseTime int64  `json:"response_time(ms),omitempty"`
	Error        string `json:"error,omitempty"`
}

type OutputAdminInfo struct {
	RequestName  string `json:"request_name"`
	RequestCount uint64 `json:"request_count,omitempty"`
	Error        string `json:"error,omitempty"`
}

var (
	siteNotLoad        = errors.New("site data not loaded. Please wait")
	siteNotConnect     = errors.New("site is unavailable")
	siteNameNotCorrect = errors.New("incorrect site requested")
	notCountRequested  = errors.New("not requested this site")
	codeNotCorrectSite = -1
	codeNotLoad        = 0
	codeOK             = 200
)

type Repository interface {
	ReadSiteInfo(s string) model.SiteResponseInfo
	ReadMinResponseSite() model.SiteResponseInfo
	ReadMaxResponseSite() model.SiteResponseInfo
	ReadCountSiteRequest(s string) (uint64, error)
	ReadCountMaxRequest() uint64
	ReadCountMinRequest() uint64
}

func NewService(r Repository) *Service {
	logger.Debug("Register service...")
	return &Service{
		repo: r,
	}
}

func NewOutputUserInfo(m model.SiteResponseInfo) OutputUserInfo {
	logger.Debug("Create new output user info...", zap.String("site", m.SiteName))
	out := OutputUserInfo{
		RequestName: m.SiteName,
	}
	if m.SiteName == "" || m.Code == codeNotLoad {
		out.Error = fmt.Sprintf("%v", siteNotLoad)
	}
	if m.Code == codeNotCorrectSite && m.SiteName != "" {
		out.Error = fmt.Sprintf("%v", siteNameNotCorrect)
	}
	if m.Code > 0 && m.Code != codeOK {
		out.Error = fmt.Sprintf("%v", siteNotConnect)
	}
	if m.Code == codeOK {
		out.ResponseTime = m.ResponseTime
	}
	return out
}

func (s Service) GetSiteInfo(site string) OutputUserInfo {
	logger.Debug("Get site response...")
	result := s.repo.ReadSiteInfo(site)
	return NewOutputUserInfo(result)
}

func (s Service) GetSiteMinResponse() OutputUserInfo {
	logger.Debug("Get Min site response...")
	result := s.repo.ReadMinResponseSite()
	return NewOutputUserInfo(result)
}

func (s Service) GetSiteMaxResponse() OutputUserInfo {
	logger.Debug("Get Max site response...")
	result := s.repo.ReadMaxResponseSite()
	return NewOutputUserInfo(result)
}

func (s Service) GetSiteStat(site string) OutputAdminInfo {
	logger.Debug("Get site stat...", zap.String("site", site))
	result, err := s.repo.ReadCountSiteRequest(site)
	out := OutputAdminInfo{
		RequestName: site,
	}
	if err != nil {
		out.Error = fmt.Sprintf("%v", siteNameNotCorrect)
	}
	if err == nil && result != 0 {
		out.RequestCount = result
	}
	if err == nil && result == 0 {
		out.Error = fmt.Sprintf("%v", notCountRequested)
	}
	return out
}

func (s Service) GetMinStat() OutputAdminInfo {
	logger.Debug("Get Min site stat...")
	res := OutputAdminInfo{
		RequestName: "Min response Endpoint requests"}
	count := s.repo.ReadCountMinRequest()
	if count == 0 {
		res.Error = fmt.Sprintf("%v", notCountRequested)
	} else {
		res.RequestCount = count
	}
	return res
}

func (s Service) GetMaxStat() OutputAdminInfo {
	logger.Debug("Get Max site stat...")
	res := OutputAdminInfo{
		RequestName: "Max response Endpoint requests"}
	count := s.repo.ReadCountMaxRequest()
	if count == 0 {
		res.Error = fmt.Sprintf("%v", notCountRequested)
	} else {
		res.RequestCount = count
	}
	return res
}
