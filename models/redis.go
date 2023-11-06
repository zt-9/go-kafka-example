package models

import (
	"context"
	"encoding/json"
	"go-kafka-example/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

type DataType interface {
	Address | Label | Transaction
}
type MetadataDB[T DataType] struct {
	db int // redis db #number
	mu sync.RWMutex
}

var client = redis.NewClient(&redis.Options{
	Addr:     config.RedisServerAddr,
	Password: "", // no password set
	DB:       0,  // use default DB
})

func NewMetadataDB[T DataType](db int) *MetadataDB[T] {
	return &MetadataDB[T]{
		db: db,
	}
}

func (m *MetadataDB[T]) KeyExists(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ctx := context.Background()

	// select database
	if _, err := client.Do(ctx, "SELECT", m.db).Result(); err != nil {
		return false, err
	}
	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (m *MetadataDB[T]) Get(key string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res T
	ctx := context.Background()
	// select database
	if _, err := client.Do(ctx, "SELECT", m.db).Result(); err != nil {
		return res, err
	}

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return res, err
	}

	return res, err
}

func (m *MetadataDB[T]) GetAll(client *redis.Client) ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []T
	ctx := context.Background()
	// select database
	if _, err := client.Do(ctx, "SELECT", m.db).Result(); err != nil {
		return res, err
	}

	// get all keys in the db
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return res, err
	}

	for _, key := range keys {
		val, err := client.Get(ctx, key).Result()
		if err != nil {
			return res, err
		}

		var data T
		err = json.Unmarshal([]byte(val), &data)
		if err != nil {
			return res, err
		}

		res = append(res, data)
	}

	return res, nil
}

func (m *MetadataDB[T]) Update(key string, data T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx := context.Background()

	// select database
	if _, err := client.Do(ctx, "SELECT", m.db).Result(); err != nil {
		return err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = client.Set(ctx, key, dataJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *MetadataDB[T]) Delete(key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ctx := context.Background()

	// select database
	if _, err := client.Do(ctx, "SELECT", m.db).Result(); err != nil {
		return 0, err
	}

	keysDeleted, err := client.Del(ctx, key).Result() // keysDeleted: the number of key being deleted. 0 means no key is deleted (the key to delete does not exists)
	if err != nil {
		return 0, err
	}
	return keysDeleted, nil
}

var AddressDB = NewMetadataDB[Address](0)
var LabelDB = NewMetadataDB[Label](1)
var TransactionDB = NewMetadataDB[Transaction](2)
