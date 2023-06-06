package handlers

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

const (
	getSite = "/url/{site}"
	getMin  = "/min"
	getMax  = "/max"
	getStat = "/stat"
)

type Hand struct {
	log     *zap.Logger
	cfg     *config.Config
	service *service.Service
}

type Handler interface {
	Register(router *chi.Mux)
}

func (h *Hand) Register(router *chi.Mux) {
	router.Get(getSite, h.GetSite)
	router.Get(getMin, h.GetMin)
	router.Get(getMax, h.GetMax)
	router.Get(getStat, h.GetStat)
}

func NewHandler(log *zap.Logger, cfg *config.Config) *Hand {
	log.Debug("Register user handler...")
	s := service.NewService(log, cfg)
	return &Hand{log: log, cfg: cfg, service: s}
}

type RespService interface {
	GetSite(site string) (float64, error)
	GetMin() (float64, string, error)
	GetMax() (float64, string, error)
	GetStat() (float64, string, error)
}

func (h *Hand) GetSite(w http.ResponseWriter, r *http.Request) {
	site := chi.URLParam(r, "site")
	h.log.Debug(fmt.Sprintf("Get Site response %s ...", site))
	time, err := h.service.GetSite(site)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Response time for %s: for %v", site, time)))
	}
	//w.WriteHeader(200)
}

func (h *Hand) GetMin(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Min Site response...")
	time, site, err := h.service.GetMin()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Minimal time response: %v for %s", time, site)))
	}
	//w.WriteHeader(200)
}

func (h *Hand) GetMax(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Max Site response...")
	time, site, err := h.service.GetMax()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Maximum time response: %v for %s", time, site)))
	}
	//w.WriteHeader(200)
}

func (h *Hand) GetStat(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Get Stat for admin...")
	time, site, err := h.service.GetStat()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
	} else {
		w.Write([]byte(fmt.Sprintf("Count stat: %v for %s", time, site)))
	}
	//w.WriteHeader(200)
}
