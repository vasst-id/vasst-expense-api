package redis

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRedisClient(t *testing.T) Cache {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	redis, err := New(&ClientOptions{
		ServiceName: "test" + "-" + "redis",
	})
	if err != nil {
		log.Fatalf("error init redis pool %s", err.Error())
	}
	return redis
}

type valueStruct struct {
	Name string
	Age  int
}

func (vs valueStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(vs)
}

func (vs *valueStruct) UnmarshalBinary(msg []byte) error {
	err := json.Unmarshal(msg, &vs)
	vs = nil
	return err
}

func TestRedis(t *testing.T) {

	t.Parallel()
	redisClient := setupRedisClient(t)
	key := "key"
	valueString := "value"
	duration := 2 * time.Minute

	t.Run("success Ping() redis server", func(t *testing.T) {
		err := redisClient.Ping()
		assert.Nil(t, err)
	})

	t.Run("success Get() empty redis value", func(t *testing.T) {
		valueString, err := redisClient.GetString(context.Background(), "key1")
		assert.Equal(t, "", valueString)
		assert.Nil(t, err)
	})

	t.Run("success GetString() a value that has been Set()", func(t *testing.T) {
		err := redisClient.SetWithExpiration(context.Background(), key, valueString, duration)
		assert.Nil(t, err)

		valueString, err := redisClient.GetString(context.Background(), key)
		assert.Equal(t, valueString, valueString)
		assert.Nil(t, err)
	})

	t.Run("success Get() a value that has been Set()", func(t *testing.T) {
		var objResult valueStruct

		err := redisClient.Set(context.Background(), key, valueStruct{
			Name: "De'Aaron Fox",
			Age:  25,
		})
		assert.Nil(t, err)

		err = redisClient.Get(context.Background(), key, &objResult)
		assert.Nil(t, err)

		expectedObj := valueStruct{
			Name: "De'Aaron Fox",
			Age:  25,
		}

		assert.Equal(t, expectedObj, objResult)
	})

	t.Run("success Remove() an existing key", func(t *testing.T) {
		err := redisClient.Set(context.Background(), key, valueString)
		ctx := context.Background()
		assert.Nil(t, err)

		err = redisClient.Remove(ctx, key)
		assert.Nil(t, err)
	})

	t.Run("success listing existing keys with Keys()", func(t *testing.T) {
		err := redisClient.Set(context.Background(), key, valueString)
		assert.Nil(t, err)
		key2 := "key2"

		err = redisClient.Set(context.Background(), key2, valueString)
		assert.Nil(t, err)

		keylist, err := redisClient.Keys(key)
		assert.Equal(t, []string{key}, keylist)
		assert.Nil(t, err)
	})

	t.Run("success getting TTL() from redis key", func(t *testing.T) {
		err := redisClient.SetWithExpiration(context.Background(), key, valueString, duration)
		assert.Nil(t, err)

		ttl, err := redisClient.TTL(key)
		assert.Nil(t, err)
		assert.Equal(t, duration, ttl)
	})

	t.Run("success getting TTL() from redis key", func(t *testing.T) {
		err := redisClient.SetWithExpiration(context.Background(), key, valueString, duration)
		assert.Nil(t, err)

		ttl, err := redisClient.TTL(key)
		assert.Nil(t, err)
		assert.Equal(t, duration, ttl)
	})

	t.Run("success MSet() and MGet()", func(t *testing.T) {
		key2 := "key2"
		var names []interface{}
		names = append(names, "Jalen Brunson")
		names = append(names, "Shai Gilgeous-Alexander")
		err := redisClient.MSet([]string{key, key2}, names)
		assert.Nil(t, err)

		multipleValues, err := redisClient.MGet([]string{key, key2})
		assert.Equal(t, names, multipleValues)
		assert.Nil(t, err)
	})

	t.Run("success MSetWithExpiration() and MGet()", func(t *testing.T) {
		key2 := "key2"
		var names []interface{}
		names = append(names, "Derrick Rose")
		names = append(names, "Bam Adebayo")
		err := redisClient.MSetWithExpiration([]string{key, key2}, names, []time.Duration{duration, duration})
		assert.Nil(t, err)

		multipleValues, err := redisClient.MGet([]string{key, key2})
		assert.Equal(t, names, multipleValues)
		assert.Nil(t, err)
	})

	t.Run("success HSET and HGET", func(t *testing.T) {
		myKey := "mykey"

		var data valueStruct
		dataSet := valueStruct{
			Name: "De'Aaron Fox",
			Age:  25,
		}
		err := redisClient.HSet(context.Background(), myKey, "field1", &dataSet)
		assert.Nil(t, err)

		err = redisClient.HGet(context.Background(), myKey, "field1", &data)
		assert.Nil(t, err)
		assert.Equal(t, dataSet, data)
	})

	t.Run("success HSET with expiration, and get TTL", func(t *testing.T) {
		myKey := "mykey"

		dataSet := valueStruct{
			Name: "De'Aaron Fox",
			Age:  25,
		}
		err := redisClient.HSetWithExpiration(context.Background(), myKey, []any{"field1", &dataSet}, duration)
		assert.Nil(t, err)

		ttl, err := redisClient.TTL(myKey)
		assert.Nil(t, err)
		assert.Equal(t, duration, ttl)
	})
}
