package main

import (
	"context"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/handlers"
	"github.com/Snegniy/testTaskResponseApi/internal/logger"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := logger.NewLogger()
	cfg := config.NewConfig(log)
	log.Debug("Create router...")
	r := chi.NewRouter()
	h := handlers.NewHandler(log, cfg)
	h.Register(r)
	start(r, log, cfg.Server.Host, cfg.Server.Port)
}

func start(r *chi.Mux, log *zap.Logger, host, port string) {
	log.Debug("Start app server")
	srv := &http.Server{
		Addr:    fmt.Sprintf(host + ":" + port),
		Handler: r,
	}

	// Запуск веб-сервера в отдельном горутине
	go func() {
		log.Info(fmt.Sprintf("Server started on %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// Ожидание сигнала для начала завершения работы
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	// Установка тайм-аута для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer log.Sync()
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("Server Shutdown: %s", err))
	}

	log.Info("Server exiting")
}
