package handlers

import (
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	srv Services
}

type Services interface {
	GetSiteInfo(site string) service.OutputUserInfo
	GetSiteMinResponse() service.OutputUserInfo
	GetSiteMaxResponse() service.OutputUserInfo
	GetSiteStat(site string) service.OutputAdminInfo
	GetMinStat() service.OutputAdminInfo
	GetMaxStat() service.OutputAdminInfo
}

func NewHandlers(s Services) Handlers {
	logger.Debug("Register user handler...")
	return Handlers{
		srv: s,
	}
}

func (h *Handlers) GetSiteResponse(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	logger.Debug("Handler call", zap.String("site", site))
	res := h.srv.GetSiteInfo(site)
	writeJSON(w, res)
}

func (h *Handlers) GetMinSiteResponse(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler call", zap.String("min", "min"))
	res := h.srv.GetSiteMinResponse()
	writeJSON(w, res)
}

func (h *Handlers) GetMaxSiteResponse(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler call", zap.String("max", "max"))
	res := h.srv.GetSiteMaxResponse()
	writeJSON(w, res)
}

func (h *Handlers) GetSiteStat(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	logger.Debug("Handler call", zap.String("admin site stat", site))
	res := h.srv.GetSiteStat(site)
	writeJSON(w, res)
}

func (h *Handlers) GetMinStat(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler call", zap.String("admin min stat", "min"))
	res := h.srv.GetMinStat()
	writeJSON(w, res)
}

func (h *Handlers) GetMaxStat(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler call", zap.String("admin max stat", "max"))
	res := h.srv.GetMaxStat()
	writeJSON(w, res)
}
