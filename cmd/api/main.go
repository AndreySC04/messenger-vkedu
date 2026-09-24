package main

import (
	"context"
	"log"

	"messenger-vkedu/internal/config"
	"messenger-vkedu/internal/db"
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
	log.Println("LOG - postgres pool ready")
}
