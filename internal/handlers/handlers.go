package handlers

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

type Hand struct {
	log     *zap.Logger
	service *service.Service
}

func NewHandler(log *zap.Logger, s *service.Service) *Hand {
	log.Debug("Register user handler...")
	return &Hand{log, s}
}

func (h *Hand) GetSiteResponse(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get Site response %s ...", site))
	time, err := h.service.GetSite(site)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Response time for %s: for %v", site, time)))
	}
}

func (h *Hand) GetMinSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min Site response...")
	time, site, err := h.service.GetMin()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Minimal time response: %v for %s", time, site)))
	}
}

func (h *Hand) GetMaxSiteResponse(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max Site response...")
	time, site, err := h.service.GetMax()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Maximum time response: %v for %s", time, site)))
	}
}

func (h *Hand) GetRequestSitesStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Stat for admin...")
	time, site, err := h.service.GetStat()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Count stat: %v for %s", time, site)))
	}
}
