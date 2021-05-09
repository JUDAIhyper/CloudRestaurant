package tool

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var RedisStoreEngine RedisStore

type RedisStore struct {
	client *redis.Client
	prefix string
}

func InitRedisStore() *RedisStore {
	config := GetConfig().RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Db,
	})
	RedisStoreEngine = RedisStore{client: client}
	return &RedisStoreEngine
}

// Set string value
func (rs *RedisStore) Set(key string, value interface{}, expiration time.Duration) {
	err := rs.client.Set(rs.prefix+key, value, expiration).Err()
	if err != nil {
		log.Println(err)
	}
}

// set json value
func (rs *RedisStore) SetJson(key string, value map[string]interface{}, expiration time.Duration) {
	jsonByte, _ := json.Marshal(value)
	err := rs.client.Set(rs.prefix+key, string(jsonByte), expiration).Err()
	if err != nil {
		log.Println(err)
	}
}

// Get value by id
func (rs *RedisStore) Get(key string) string {
	val, err := rs.client.Get(rs.prefix+key).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	return val
}

// GetJson get json value by id
func (rs *RedisStore) GetJson(key string) map[string]interface{} {
	val, err := rs.client.Get(rs.prefix+key).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		log.Println(err.Error())
		return nil
	}
	return result
}

// Del value by id
func (rs *RedisStore) Del(key string) bool {
	err := rs.client.Del(rs.prefix+key).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
