package main

import (
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/handlers"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"github.com/Snegniy/testTaskResponseApi/internal/response"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/Snegniy/testTaskResponseApi/pkg/graceful"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/chi"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg.IsDebug)

	log.Debug("Create router...")
	router := chi.NewRouter()

	r := repository.NewRepository(log, cfg.UrlRepo.SitesFile)
	s := service.NewService(log, r)
	h := handlers.NewHandlers(log, s)

	Register(router, h)

	graceful.StartServer(router, log, cfg.Server.Host, cfg.Server.Port)
	response.Response(r, cfg.UrlRepo.Timeout, cfg.UrlRepo.Refresh)
}

func Register(router *chi.Mux, h handlers.Handlers) {
	router.Get("/url/{site}", h.GetSiteResponse)
	router.Get("/min", h.GetMinSiteResponse)
	router.Get("/max", h.GetMaxSiteResponse)
	router.Get("/stat/url/{site}", h.GetSiteStat)
	router.Get("/stat/min", h.GetRequestMinSitesStat)
	router.Get("/stat/max", h.GetRequestMaxSitesStat)
}
