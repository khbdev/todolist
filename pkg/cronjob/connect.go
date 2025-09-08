
package cronjob

import (
    "github.com/hibiken/asynq"
    "github.com/redis/go-redis/v9"
    "context"
    "log"
)

var (
    RedisClient *redis.Client
    AsynqClient *asynq.Client
)

func ConnectRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })

    _, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal("Redisga ulanishda xato:", err)
    }

    AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
        Addr: "localhost:6379",
    })
}
