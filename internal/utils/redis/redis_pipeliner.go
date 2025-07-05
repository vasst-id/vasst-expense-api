package redis

import (
	"context"
	redis "github.com/go-redis/redis/v8"
)

type Pipeliner interface {
	redis.StatefulCmdable
	Len() int
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Process(ctx context.Context, cmd redis.Cmder) error
	Close() error
	Discard() error
	Exec(ctx context.Context) ([]redis.Cmder, error)
}
