package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/adapter/controllers"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/adapter/postgres"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/adapter/postgres/repositories"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/config"
	"github.com/tulioguaraldob/luizalabs-logistica-challenge/internal/domain/services"
)

func main() {
	// Config
	if err := config.LoadEnvs(); err != nil {
		log.Fatalf("Failed to load environment variables. Details: %s", err.Error())
	}

	// Postgres
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Env.PostgresUser,
		config.Env.PostgresPassword,
		config.Env.PostgresHost,
		config.Env.PostgresPort,
		config.Env.PostgresDb,
	)
	db, err := postgres.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close(db)

	// Repositories
	pr := repositories.NewProductRepository(db)
	or := repositories.NewOrderRepository(db)
	ur := repositories.NewUserRepository(db)
	opr := repositories.NewOrderProductRepository(db)

	// Services
	us := services.NewUserService(ur, or, pr, opr)
	ors := services.NewOrderService(or, opr, ur)

	// Controllers
	uc := controllers.NewUserController(us)
	oc := controllers.NewOrderController(ors)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("I heard that Tulio was approved on LuizaLabs! :)"))
	})
	mux.HandleFunc("GET /user/{id}", uc.Get)
	mux.HandleFunc("POST /user/upload", uc.PostUsersData)
	mux.HandleFunc("GET /order/{id}", oc.GetByID)
	mux.HandleFunc("GET /orders", oc.Get)

	port := config.Env.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func() {
		log.Printf("Running on port: %s\n", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve HTTP server. Details: %s", err.Error())
		}
	}()
	<-ctx.Done()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server. Details: %s", err.Error())
	}
}
