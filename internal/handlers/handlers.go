package handlers

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
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
	GetSiteInfo(site string) service.OutputUserInfo
	GetSiteMinResponse() service.OutputUserInfo
	GetSiteMaxResponse() service.OutputUserInfo
	GetSiteStat(site string) service.OutputAdminInfo
	GetMinStat() service.OutputAdminInfo
	GetMaxStat() service.OutputAdminInfo
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
	res := h.srv.GetSiteInfo(site)
	writeJSON(w, res)
}

func (h *Handlers) GetMinSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min site response...")
	res := h.srv.GetSiteMinResponse()
	writeJSON(w, res)
}

func (h *Handlers) GetMaxSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max site response...")
	res := h.srv.GetSiteMaxResponse()
	writeJSON(w, res)
}

func (h *Handlers) GetSiteStat(w http.ResponseWriter, r *http.Request) {
	//_, claims, _ := jwtauth.FromContext(r.Context())
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get Stat count requests site %s for admin...", site))
	res := h.srv.GetSiteStat(site)
	writeJSON(w, res)
}

func (h *Handlers) GetMinStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min response count stat for admin...")
	res := h.srv.GetMinStat()
	writeJSON(w, res)
}

func (h *Handlers) GetMaxStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max response count stat for admin...")
	res := h.srv.GetMinStat()
	writeJSON(w, res)
}
