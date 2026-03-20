package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"salary-calculator/internal/config"
	"salary-calculator/internal/middleware"
	"salary-calculator/internal/repository"
	"syscall"
	"time"
)

func Run() error {

	cfg := config.Load()

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		return err
	}

	defer db.Close()

	router := newRouter()

	handler := middleware.Logging(router)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	return startServer(server)
}

func newRouter() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}

func startServer(server *http.Server) error {

	go func() {

		log.Println("server started on", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}

	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	return server.Shutdown(ctx)
}
