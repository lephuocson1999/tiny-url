package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"

	"tiny-url/internal/api"
	"tiny-url/internal/app"
	"tiny-url/internal/config"
	"tiny-url/internal/repo"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Postgres setup
	db, err := sql.Open("pgx", dsn(cfg))
	if err != nil {
		log.Printf("failed to connect to Postgres: %v, using in-memory repo", err)
	}
	var repository app.Repository
	if db != nil && db.Ping() == nil {
		repository = app.NewPostgresRepo(db)
	} else {
		repository = repo.NewMemoryRepo()
	}

	// Redis and SafeIDGenerator setup
	var cache app.Cache
	var idGen app.IDGenerator
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host + ":" + itoa(cfg.Redis.Port),
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis for ID generation: %v", err)
	} else {
		cache = app.NewRedisCache(rdb)
		idGen = app.NewSafeIDGenerator(rdb, db, "url:id:counter", 1000)
		if err := idGen.(*app.SafeIDGenerator).RestoreFromDB(context.Background()); err != nil {
			log.Fatalf("failed to restore ID counter from DB: %v", err)
		}
	}

	service := app.NewURLShortenerService(repository, cache, idGen)
	handler := api.NewHandler(service)

	r := gin.Default()
	r.POST("/api/v1/shorten", handler.ShortenURL)
	r.GET("/api/v1/:shortCode", handler.Redirect)
	r.GET("/api/v1/stats/:shortCode", handler.GetStats)

	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", cfg.Server.Port)
	}
	log.Printf("Server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func dsn(cfg *config.Config) string {
	return "host=" + cfg.Database.Host +
		" port=" + itoa(cfg.Database.Port) +
		" user=" + cfg.Database.Username +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.Name +
		" sslmode=disable"
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
