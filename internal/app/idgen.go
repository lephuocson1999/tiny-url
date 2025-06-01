package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SafeIDGenerator struct {
	rdb         *redis.Client
	db          *sql.DB
	key         string
	backupEvery int64
}

func NewSafeIDGenerator(rdb *redis.Client, db *sql.DB, key string, backupEvery int64) *SafeIDGenerator {
	return &SafeIDGenerator{rdb: rdb, db: db, key: key, backupEvery: backupEvery}
}

func (g *SafeIDGenerator) NextID(ctx context.Context) (int64, error) {
	id, err := g.rdb.Incr(ctx, g.key).Result()
	if err != nil {
		return 0, err
	}

	fmt.Println("id", id)
	go func(id int64) {
		_, _ = g.db.ExecContext(context.Background(),
			"UPDATE id_counters SET value = $1 WHERE name = 'url_id'", id)
	}(id)

	return id, nil
}

func (g *SafeIDGenerator) RestoreFromDB(ctx context.Context) error {
	var dbVal int64
	err := g.db.QueryRowContext(ctx, "SELECT value FROM id_counters WHERE name = 'url_id'").Scan(&dbVal)
	if err != nil {
		return err
	}
	redisVal, _ := g.rdb.Get(ctx, g.key).Int64()
	if dbVal != redisVal {
		return g.rdb.Set(ctx, g.key, dbVal, 0).Err()
	}
	return nil
}
