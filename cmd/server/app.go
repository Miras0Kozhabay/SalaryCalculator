package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"salary-calculator/internal/calculator"
	"salary-calculator/internal/config"
	"salary-calculator/internal/handlers"
	"salary-calculator/internal/middleware"
	"salary-calculator/internal/repository"
	"salary-calculator/internal/services"
)

func Run() error {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// Build PostgreSQL DSN with SSL mode configuration
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Configure connection pool for production
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Printf("✓ Connected to database: %s@%s:%d/%s",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	calc := calculator.NewCalculator(cfg.MCI)
	repo := repository.NewPostgresRepository(db)
	salaryService := services.NewSalaryService(calc, repo)
	salaryHandler := handlers.NewSalaryHandler(salaryService)

	router := http.NewServeMux()
	router.HandleFunc("POST /api/calculate", salaryHandler.Calculate)
	router.HandleFunc("GET /api/history", salaryHandler.History)
	router.HandleFunc("GET /api/mci", salaryHandler.MCI)
	router.Handle("/", http.FileServer(http.Dir("./web")))

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler := middleware.Logging(router)

	addr := ":" + strconv.Itoa(cfg.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return startServer(server, addr)
}

func startServer(server *http.Server, addr string) error {
	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("⏹️  Graceful shutdown initiated...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("✓ Server stopped gracefully")
	return nil
}
