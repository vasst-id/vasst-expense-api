package healthcheck

import (
	"context"
	"database/sql"
	"time"

	logs "github.com/vasst-id/vasst-expense-api/internal/utils/logger"
	// "github.com/julofinance/julo-go-library/redis"
)

// Config used to add information when constructing component
type Config struct {
	Name    string
	Timeout time.Duration
}

type Option func(*healthService)

func WithComponent(c Component) Option {
	return func(h *healthService) {
		h.register(c)
	}
}

func WithDB(db *sql.DB, cfg Config) Option {
	return func(h *healthService) {
		c := Component{
			Name:    cfg.Name,
			Timeout: cfg.Timeout,
			CheckFunc: func(ctx context.Context) error {
				return db.PingContext(ctx)
			},
		}
		h.register(c)
	}
}

// func WithRedis(cache redis.Cache, cfg Config) Option {
// 	return func(h *healthService) {
// 		c := Component{
// 			Name:    cfg.Name,
// 			Timeout: cfg.Timeout,
// 			CheckFunc: func(ctx context.Context) error {
// 				return cache.Ping()
// 			},
// 		}
// 		h.register(c)
// 	}
// }

func WithLogger(logger *logs.Logger) Option {
	return func(h *healthService) {
		h.logger = logger
	}
}
