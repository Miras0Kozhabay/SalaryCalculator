package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"salary-calculator/internal/calculator"
	"salary-calculator/internal/config"
	"salary-calculator/internal/handlers"
	"salary-calculator/internal/middleware"
	"salary-calculator/internal/repository"
	"salary-calculator/internal/services"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func Run() error {
	// Load config
	cfg := config.Load()

	// Connect to PostgreSQL
	dsn := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Init calculator
	calc := calculator.NewCalculator(cfg.MCI)

	// Init repository
	repo := repository.NewPostgresRepository(db)

	// Init service
	salaryService := services.NewSalaryService(calc, repo)

	// Init handlers
	salaryHandler := handlers.NewSalaryHandler(salaryService)

	// Router
	router := http.NewServeMux()
	router.HandleFunc("/calculate", salaryHandler.Calculate)
	router.HandleFunc("/history", salaryHandler.History)
	router.HandleFunc("/mci", salaryHandler.MCI)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Apply middleware
	handler := middleware.Logging(router)

	// HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	return startServer(server)
}

func startServer(server *http.Server) error {
	// Start server in goroutine
	go func() {
		log.Println("server started on", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
