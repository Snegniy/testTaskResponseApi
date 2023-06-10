package service

import (
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"go.uber.org/zap"
	"strconv"
)

type Service struct {
	repo Repository
}

type OutputUserInfo struct {
	RequestName  string `json:"request_name"`
	ResponseTime string `json:"response_time"`
}

type OutputAdminInfo struct {
	RequestName  string `json:"request_name"`
	RequestCount string `json:"request_count"`
}

var (
	siteNotLoad        = errors.New("error: site data not loaded. Please wait")
	siteNotConnect     = errors.New("error: site is unavailable")
	siteNameNotCorrect = errors.New("error: incorrect site requested")
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
	if m.Code == codeNotCorrectSite {
		out.ResponseTime = fmt.Sprintf("%v", siteNameNotCorrect)
	}
	if m.Code == codeNotLoad {
		out.ResponseTime = fmt.Sprintf("%v", siteNotLoad)
	}
	if m.Code > 0 && m.Code != codeOK {
		out.ResponseTime = fmt.Sprintf("%v", siteNotConnect)
	}
	if m.Code == codeOK {
		t := strconv.Itoa(int(m.ResponseTime))
		out.ResponseTime = fmt.Sprintf("%sms", t)
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
	c := strconv.Itoa(int(result))
	out := OutputAdminInfo{
		RequestName:  site,
		RequestCount: c,
	}
	if err != nil {
		out.RequestCount = fmt.Sprintf("%v", siteNameNotCorrect)
	}
	return out
}

func (s Service) GetMinStat() OutputAdminInfo {
	logger.Debug("Get Min site stat...")
	return OutputAdminInfo{
		RequestName:  "Min response Endpoint requests",
		RequestCount: strconv.Itoa(int(s.repo.ReadCountMinRequest())),
	}
}

func (s Service) GetMaxStat() OutputAdminInfo {
	logger.Debug("Get Max site stat...")
	return OutputAdminInfo{
		RequestName:  "Max response Endpoint requests",
		RequestCount: strconv.Itoa(int(s.repo.ReadCountMaxRequest())),
	}
}
