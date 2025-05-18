package main

import (
	"log"

	"github.com/kuluruvineeth/social-go/internal/db"
	"github.com/kuluruvineeth/social-go/internal/env"
	"github.com/kuluruvineeth/social-go/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social_go?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdletime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdletime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("Connected to the database")

	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))

}
