package redis

import (
	"context"
	"encoding"
	"fmt"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

//go:generate mockgen -source=redis_universal.go -package=redis -destination=redis_universal_mock.go
type (
	ClientOptions struct {
		ServiceName     string
		Address         []string
		DB              int
		Username        string
		Password        string
		PoolSize        int
		MinIdleConns    int
		ConnMaxLifetime time.Duration
		ConnMaxIdleTime time.Duration
		PoolTimeout     time.Duration
		DialTimeout     time.Duration
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		DataDogTracer   bool //trace to be on or off
	}

	//go:generate mockgen -destination=redis_universal_test.go -package=redis . Cache
	Cache interface {
		TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) error

		SetWithExpiration(ctx context.Context, key string, value interface{}, duration time.Duration) error
		Set(ctx context.Context, key string, value interface{}) error
		SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
		Get(ctx context.Context, key string, value interface{}) error
		GetString(ctx context.Context, key string) (string, error)
		GetInt(ctx context.Context, key string) (int, error)
		Incr(ctx context.Context, key string) (int64, error)
		IncrBy(ctx context.Context, key string, value int64) (int64, error)

		HGet(ctx context.Context, key, field string, data any) error
		HSet(ctx context.Context, key string, values ...any) error
		HSetWithExpiration(ctx context.Context, key string, values []any, ttl time.Duration) error

		MSet(keys []string, values []interface{}) error
		MGet(key []string) ([]interface{}, error)
		MSetWithExpiration(keys []string, values []interface{}, ttls []time.Duration) error

		LPush(ctx context.Context, key string, value interface{}) (int64, error)

		Keys(string) ([]string, error)
		TTL(key string) (time.Duration, error)

		Remove(ctx context.Context, key string) error
		Ping() error
	}

	redisUniversalClient struct {
		r redis.UniversalClient
	}
)

var ctx = context.Background()

func New(co *ClientOptions) (Cache, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        co.Address,
		DB:           co.DB,
		Username:     co.Username,
		Password:     co.Password,
		PoolSize:     co.PoolSize,
		MinIdleConns: co.MinIdleConns,
		MaxConnAge:   co.ConnMaxLifetime,
		IdleTimeout:  co.ConnMaxIdleTime,
		PoolTimeout:  co.PoolTimeout,
		DialTimeout:  co.DialTimeout,
		ReadTimeout:  co.ReadTimeout,
		WriteTimeout: co.WriteTimeout,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, errors.New("failed to connect to redis")
	}
	if co.DataDogTracer {
		redistrace.WrapClient(client, redistrace.WithServiceName(co.ServiceName))
	}

	return &redisUniversalClient{r: client}, nil
}

func (c *redisUniversalClient) HSet(ctx context.Context, key string, values ...any) error {

	if err := check(c); err != nil {
		return err
	}

	_, err := c.r.HSet(ctx, key, values...).Result()

	return err
}

func (c *redisUniversalClient) HSetWithExpiration(ctx context.Context, key string, values []any, ttl time.Duration) error {

	if err := check(c); err != nil {
		return err
	}

	_, err := c.r.HSet(ctx, key, values...).Result()
	if err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}

	_, err = c.r.Expire(ctx, key, ttl).Result()
	if err != nil {
		c.r.Del(ctx, key)
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}

	return nil
}

func (c *redisUniversalClient) HGet(ctx context.Context, key, field string, data any) error {

	if _, ok := data.(encoding.BinaryUnmarshaler); !ok {
		return errors.New(fmt.Sprintf("failed to get cache with key %s!: redis: can't unmarshal (implement encoding.BinaryUnmarshaler)", key))
	}

	if err := check(c); err != nil {
		return err
	}

	val, err := c.r.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return errors.Wrapf(err, "key %s does not exists", key)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to get key %s!", key)
	}

	if err := data.(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (c *redisUniversalClient) TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
	if err := check(c); err != nil {
		return err
	}

	_, err := c.r.TxPipelined(ctx, fn)
	return err
}

func (c *redisUniversalClient) Ping() error {
	if _, err := c.r.Ping(ctx).Result(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *redisUniversalClient) SetWithExpiration(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Set(ctx, key, value, duration).Result(); err != nil {
		return errors.Wrapf(err, "failed to set cache with key %s!", key)
	}
	return nil
}

func (c *redisUniversalClient) Set(ctx context.Context, key string, value interface{}) error {
	if err := check(c); err != nil {
		return err
	}

	return c.SetWithExpiration(ctx, key, value, 0)
}

func (c *redisUniversalClient) GetString(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}

	if err := check(c); err != nil {
		return "", err
	}

	val, err := c.r.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		return "", errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) Get(ctx context.Context, key string, data interface{}) error {
	if _, ok := data.(encoding.BinaryUnmarshaler); !ok {
		return errors.New(fmt.Sprintf("failed to get cache with key %s!: redis: can't unmarshal (implement encoding.BinaryUnmarshaler)", key))
	}

	if err := check(c); err != nil {
		return err
	}

	val, err := c.r.Get(ctx, key).Result()

	if err == redis.Nil {
		return errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to get key %s!", key)
	}

	if err := data.(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(val)); err != nil {
		return err
	}

	return nil
}

func (c *redisUniversalClient) Keys(pattern string) ([]string, error) {
	if err := check(c); err != nil {
		return []string{}, err
	}

	return c.r.Keys(ctx, pattern).Result()
}

func (c *redisUniversalClient) Remove(ctx context.Context, key string) error {
	if err := check(c); err != nil {
		return err
	}

	if _, err := c.r.Del(ctx, key).Result(); err != nil {
		return errors.Wrapf(err, "failed to remove key %s!", key)
	}

	return nil
}

func check(c *redisUniversalClient) error {
	if c.r == nil {
		return errors.New("redis client is not connected")
	}

	return nil
}

func (c *redisUniversalClient) MGet(key []string) ([]interface{}, error) {
	if err := check(c); err != nil {
		return nil, err
	}

	val, err := c.r.MGet(ctx, key...).Result()
	if err == redis.Nil {
		return nil, errors.Wrapf(err, "key %s does not exits", key)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key %s!", key)
	}

	return val, nil
}

func (c *redisUniversalClient) MSetWithExpiration(keys []string, values []interface{}, ttls []time.Duration) error {
	if len(ttls) < len(keys) {
		return errors.New("values or ttls must have same count with keys")
	}

	var failedKeys []string

	if err := c.MSet(keys, values); err != nil {
		return errors.WithStack(err)
	}

	for i := range keys {
		if _, err := c.r.Expire(ctx, keys[i], ttls[i]).Result(); err != nil {
			failedKeys = append(failedKeys, keys[i])
			c.r.Del(ctx, keys[i])
		}
	}

	if len(failedKeys) > 0 {
		return errors.New("failed to set some keys " + strings.Join(failedKeys, ","))
	}
	return nil
}

func (c *redisUniversalClient) MSet(keys []string, values []interface{}) error {
	if err := check(c); err != nil {
		return errors.WithStack(err)
	}

	if len(values) < len(keys) {
		return errors.New("values or ttls must have same count with keys")
	}

	var pairs []interface{}

	for i := range keys {
		pairs = append(pairs, keys[i], values[i])
	}
	_, err := c.r.MSet(ctx, pairs...).Result()

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *redisUniversalClient) TTL(key string) (duration time.Duration, err error) {
	if err = check(c); err != nil {
		return duration, err
	}

	duration, err = c.r.TTL(ctx, key).Result()

	if err != nil {
		return duration, errors.Wrapf(err, "failed to get TTL with key %s!", key)
	}

	return duration, nil
}

func (c *redisUniversalClient) LPush(ctx context.Context, key string, value interface{}) (listLength int64, err error) {
	if err = check(c); err != nil {
		return 0, err
	}

	listLength, err = c.r.LPush(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}

	return listLength, nil
}

func (c *redisUniversalClient) GetInt(ctx context.Context, key string) (value int, err error) {
	if err = check(c); err != nil {
		return 0, err
	}

	val, err := c.r.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (c *redisUniversalClient) Incr(ctx context.Context, key string) (value int64, err error) {
	if err = check(c); err != nil {
		return 0, err
	}

	value, err = c.r.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (c *redisUniversalClient) IncrBy(ctx context.Context, key string, value int64) (result int64, err error) {
	if err = check(c); err != nil {
		return 0, err
	}

	result, err = c.r.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *redisUniversalClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	if err := check(c); err != nil {
		return false, err
	}

	return c.r.SetNX(ctx, key, value, expiration).Result()
}
