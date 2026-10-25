package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messenger-vkedu/internal/auth"
	"messenger-vkedu/internal/config"
	"messenger-vkedu/internal/db"
	"messenger-vkedu/internal/handler"
	"messenger-vkedu/internal/middleware"
	"messenger-vkedu/internal/user"

	"github.com/gin-gonic/gin"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg, err := config.ConfigLoad()
	if err != nil {
		log.Fatalf("LOG - load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	postgresPool, err := db.NewPostgresPool(ctx, cfg.DB.URL())
	if err != nil {
		log.Fatalf("LOG - init db: %v", err)
	}
	defer postgresPool.Close()

	userRepo := user.NewRepository(postgresPool)
	refreshTokenRepo := auth.NewRepository(postgresPool)
	tokenManager := auth.NewTokenManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	authService := auth.NewService(userRepo, refreshTokenRepo, tokenManager)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	authRoutes := apiV1.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/refresh", authHandler.Refresh)
	authRoutes.POST("/logout", authHandler.Logout)

	protectedRoutes := apiV1.Group("")
	protectedRoutes.Use(middleware.RequireAuth(tokenManager))
	protectedRoutes.GET("/me", authHandler.Me)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf("LOG - server started on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("LOG - run server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("LOG - shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("LOG - server shutdown: %v", err)
	}
}
