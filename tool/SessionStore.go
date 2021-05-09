package tool

import "time"

func SetSession(key string, value map[string]interface{}) {
	redisStore := RedisStoreEngine
	redisStore.SetJson(key, value, time.Minute*60*2)
}

func GetSession(key string) map[string]interface{} {
	redisStore := RedisStoreEngine
	return redisStore.GetJson(key)
}
