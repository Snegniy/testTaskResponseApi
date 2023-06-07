package handlers

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	log     *zap.Logger
	service *service.Service
	srv     Services
}

type outputStat struct {
	SiteName string `json:"site_name"`
	Count    uint64 `json:"count_requests"`
}

type outputError struct {
	SiteName string `json:"site_name"`
	Error    string `json:"error"`
}

type Services interface {
	GetSiteInfo(site string) (model.SiteResponseInfo, error)
	GetSiteMinResponse() (model.SiteResponseInfo, error)
	GetSiteMaxResponse() (model.SiteResponseInfo, error)
	GetSiteStat(site string) (uint64, error)
	GetMinStat() uint64
	GetMaxStat() uint64
}

func NewHandlers(log *zap.Logger, s *service.Service) Handlers {
	log.Debug("Register user handler...")
	return Handlers{
		log:     log,
		service: s,
	}
}

func (h *Handlers) GetSiteResponse(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get site %s response...", site))
	res, err := h.service.GetSiteInfo(site)
	if err != nil {
		writeErrorJSON(w, outputError{site, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
	h.log.Debug(fmt.Sprintf("Get site %s response - OK!!", site))
}

func (h *Handlers) GetMinSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min site response...")
	res, err := h.service.GetSiteMinResponse()
	if err != nil {
		writeErrorJSON(w, outputError{res.SiteName, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
	h.log.Debug("Get Min site response - OK!")
}

func (h *Handlers) GetMaxSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max site response...")
	res, err := h.service.GetSiteMaxResponse()
	if err != nil {
		writeErrorJSON(w, outputError{res.SiteName, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
	h.log.Debug("Get Max site response - OK!")
}

func (h *Handlers) GetSiteStat(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get Stat count requests site %s for admin...", site))
	res, err := h.service.GetSiteStat(site)
	if err != nil {
		writeErrorJSON(w, outputError{site, fmt.Sprintf("%s", err)})
	} else {
		writeStatJSON(w, outputStat{site, res})
	}
	h.log.Debug(fmt.Sprintf("Get Stat count requests site %s for admin - OK!", site))
}

func (h *Handlers) GetMinStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min response count stat for admin...")
	res := h.service.GetMinStat()
	writeStatJSON(w, outputStat{"Min response request", res})
	h.log.Debug("Get Min response count stat for admin - OK!")
}

func (h *Handlers) GetMaxStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max response count stat for admin...")
	res := h.service.GetMinStat()
	writeStatJSON(w, outputStat{"Max response request", res})
	h.log.Debug("Get Max response count stat for admin - OK!")
}
