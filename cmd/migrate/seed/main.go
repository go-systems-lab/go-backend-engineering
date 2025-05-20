package main

import (
	"log"

	"github.com/kuluruvineeth/social-go/internal/db"
	"github.com/kuluruvineeth/social-go/internal/env"
	"github.com/kuluruvineeth/social-go/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social_go?sslmode=disable")

	conn, err := db.NewDB(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
