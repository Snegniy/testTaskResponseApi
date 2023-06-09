package handlers

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	log *zap.Logger
	srv Services
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

func NewHandlers(log *zap.Logger, s Services) Handlers {
	log.Debug("Register user handler...")
	return Handlers{
		log: log,
		srv: s,
	}
}

func (h *Handlers) GetSiteResponse(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get site %s response...", site))
	res, err := h.srv.GetSiteInfo(site)
	if err != nil {
		writeErrorJSON(w, outputError{site, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
}

func (h *Handlers) GetMinSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min site response...")
	res, err := h.srv.GetSiteMinResponse()
	if err != nil {
		writeErrorJSON(w, outputError{res.SiteName, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
}

func (h *Handlers) GetMaxSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max site response...")
	res, err := h.srv.GetSiteMaxResponse()
	if err != nil {
		writeErrorJSON(w, outputError{res.SiteName, fmt.Sprintf("%s", err)})
	} else {
		writeInfoJSON(w, res)
	}
}

func (h *Handlers) GetSiteStat(w http.ResponseWriter, r *http.Request) {
	//_, claims, _ := jwtauth.FromContext(r.Context())
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get Stat count requests site %s for admin...", site))
	res, err := h.srv.GetSiteStat(site)
	if err != nil {
		writeErrorJSON(w, outputError{site, fmt.Sprintf("%s", err)})
	} else {
		writeStatJSON(w, outputStat{site, res})
	}
}

func (h *Handlers) GetMinStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min response count stat for admin...")
	res := h.srv.GetMinStat()
	writeStatJSON(w, outputStat{"Min response request", res})
	h.log.Debug("Get Min response count stat for admin - OK!")
}

func (h *Handlers) GetMaxStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max response count stat for admin...")
	res := h.srv.GetMinStat()
	writeStatJSON(w, outputStat{"Max response request", res})
	h.log.Debug("Get Max response count stat for admin - OK!")
}
