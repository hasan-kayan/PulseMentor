package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/app"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app new: %v", err)
	}

	srv := &http.Server{
		Addr:         cfg.ServerAddr(),
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.ServerAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	_ = a.Close(ctx)
	log.Println("shutdown complete")
}
