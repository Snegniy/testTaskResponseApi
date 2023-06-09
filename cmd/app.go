package main

import (
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/handlers"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"github.com/Snegniy/testTaskResponseApi/internal/response"
	"github.com/Snegniy/testTaskResponseApi/internal/service"
	"github.com/Snegniy/testTaskResponseApi/pkg/graceful"
	"github.com/Snegniy/testTaskResponseApi/pkg/jwt"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg.ModeWork.IsDebug)

	log.Debug("Create router...")
	router := chi.NewRouter()

	r := repository.NewRepository(log, cfg.UrlRepo.SitesFile)
	s := service.NewService(log, r)
	h := handlers.NewHandlers(log, s)

	tokenAuth := jwt.NewJWT()

	Register(router, h, *tokenAuth, cfg.ModeWork.AuthAdmin)
	go cronjob.Response(r, cfg.UrlRepo.Timeout, cfg.UrlRepo.Refresh, log)
	graceful.StartServer(router, log, cfg.Server.HostPort)
}

func Register(router *chi.Mux, h handlers.Handlers, t jwtauth.JWTAuth, auth string) {
	router.Get("/url/{site}", h.GetSiteResponse)
	router.Get("/min", h.GetMinSiteResponse)
	router.Get("/max", h.GetMaxSiteResponse)

	router.Group(func(router chi.Router) {
		if auth == "jwt" {
			router.Use(jwtauth.Verifier(&t))
			router.Use(jwtauth.Authenticator)
		}
		router.Get("/stat/url/{site}", h.GetSiteStat)
		router.Get("/stat/min", h.GetMinStat)
		router.Get("/stat/max", h.GetMaxStat)
	})
}
