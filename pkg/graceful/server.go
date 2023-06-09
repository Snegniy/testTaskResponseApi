package graceful

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(r *chi.Mux, log *zap.Logger, host string) {
	log.Debug("Start app server")

	srv := &http.Server{
		Addr:    host,
		Handler: r,
	}

	go func() {
		log.Info(fmt.Sprintf("Server started on %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("Server Shutdown: %s", err))
	}

	log.Info("Server exiting")
}
