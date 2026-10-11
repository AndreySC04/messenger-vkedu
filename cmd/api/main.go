package main

import (
	"context"
	"log"

	"messenger-vkedu/internal/auth"
	"messenger-vkedu/internal/config"
	"messenger-vkedu/internal/db"
	"messenger-vkedu/internal/user"
)

func main() {
	cfg, err := config.ConfigLoad()
	if err != nil {
		log.Fatalf("LOG - load config: %v", err)
	}

	connURL := cfg.DB.URL()
	ctx := context.Background()
	postgresPool, err := db.NewPostgresPool(ctx, connURL)
	if err != nil {
		log.Fatalf("LOG - init db: %v", err)
	}
	defer postgresPool.Close()

	userRepo := user.NewRepository(postgresPool)
	refreshTokenRepo := auth.NewRepository(postgresPool)
	tokenManager := auth.NewTokenManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	authService := auth.NewService(userRepo, refreshTokenRepo, tokenManager)

	_ = authService
	log.Println("LOG - service started")
}
