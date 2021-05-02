package db

import (
	"app-auth/config"
	"fmt"
	"os"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/go-redis/redis"
)

func ConnectRedis() *redis.Client {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,  // no password set
		DB:       config.RedisDB, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pong)

	return client
}

var client = ConnectRedis()

func Set(key string, value string) error {
	duration := time.Duration(24 * time.Hour)

	return client.Set(key, value, duration).Err()
}

func SetObject(key string, data map[string]interface{}) error {
	return client.HMSet(key, data).Err()
}

func Get(key string) string {
	res, err := client.Get(key).Result()
	if err != nil {
		log.Print(err)
	}

	return res
}

// https://id.scaratec.com/user/user-invite-confirm/62a96288-d95e-4b1a-8d38-3b535014a5d9
func GetObject(key string) ([]interface{}, error) {
	return client.HMGet(key, "app", "user_id", "new_user", "team_id", "user_email", "organisation_id", "signup_url", "app_redirect_url").Result()
}

func GetObjectByKey(key string) (map[string]string, error) {
	return client.HGetAll(key).Result()
}

func Del(key string) int64 {
	res, err := client.Del(key).Result()
	if err != nil {
		log.Print(err)
	}

	return res
}
