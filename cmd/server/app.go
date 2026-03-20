package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	cfg := config.Load()

	dsn := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}
	log.Println("connected to database")

	calc := calculator.NewCalculator(cfg.MCI)
	repo := repository.NewPostgresRepository(db)
	salaryService := services.NewSalaryService(calc, repo)
	salaryHandler := handlers.NewSalaryHandler(salaryService)

	router := http.NewServeMux()
	router.HandleFunc("/api/calculate", salaryHandler.Calculate)
	router.HandleFunc("/api/history", salaryHandler.History)
	router.HandleFunc("/api/mci", salaryHandler.MCI)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler := middleware.Logging(router)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return startServer(server)
}

func startServer(server *http.Server) error {
	go func() {
		log.Printf("server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
