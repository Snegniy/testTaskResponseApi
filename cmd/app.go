package main

import (
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/cronjob"
	"github.com/Snegniy/testTaskResponseApi/internal/handlers"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/Snegniy/testTaskResponseApi/pkg/graceful"
	"github.com/Snegniy/testTaskResponseApi/pkg/jwt"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	cfg := config.NewConfig()
	logger.Init(cfg.ModeWork.IsDebug)

	logger.Debug("Create router...")
	router := chi.NewRouter()

	r := repository.NewRepository(cfg.UrlRepo.SitesFile)
	s := service.NewService(r)
	h := handlers.NewHandlers(s)

	Register(router, h, cfg.ModeWork.AuthAdmin)
	go cronjob.SiteCheckResponse(r, cfg.UrlRepo.Timeout, cfg.UrlRepo.Refresh)
	graceful.StartServer(router, cfg.Server.HostPort)
}

func Register(router *chi.Mux, h handlers.Handlers, auth string) {
	router.Get("/min", h.GetMinSiteResponse)
	router.Get("/max", h.GetMaxSiteResponse)
	router.Get("/{site}", h.GetSiteResponse)

	router.Group(func(router chi.Router) {
		if auth == "jwt" {
			tokenAuth := jwt.NewJWT()
			router.Use(jwtauth.Verifier(tokenAuth))
			router.Use(jwtauth.Authenticator)
		}
		router.Get("/stat/min", h.GetMinStat)
		router.Get("/stat/max", h.GetMaxStat)
		router.Get("/stat/{site}", h.GetSiteStat)
	})
}
