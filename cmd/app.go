package main

import (
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/handlers"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"github.com/Snegniy/testTaskResponseApi/internal/responser"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/Snegniy/testTaskResponseApi/pkg/graceful"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/chi"
)

func main() {
	log := logger.NewLogger()
	cfg := config.NewConfig(log)
	log.Debug("Create router...")
	router := chi.NewRouter()
	r := repository.NewRepository(log, cfg.UrlRepo.SitesFile)
	s := service.NewService(log, r)
	h := handlers.NewHandler(log, s)
	Register(router, h)
	graceful.StartServer(router, log, cfg.Server.Host, cfg.Server.Port)
	responser.Response(r, cfg.UrlRepo.Timeout, cfg.UrlRepo.Refresh, log)
}

func Register(router *chi.Mux, h *handlers.Hand) {
	router.Get("/url/{site}", h.GetSiteResponse)
	router.Get("/min", h.GetMinSiteResponse)
	router.Get("/max", h.GetMaxSiteResponse)
	router.Get("/stat", h.GetRequestSitesStat)
}
