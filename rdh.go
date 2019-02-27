package rdh

import (
    "strings"
    "time"
    "github.com/go-redis/redis"
)

type redisHandle struct {
    IsCluster       bool
    ClusterClient   *redis.ClusterClient
    OneNodeClient   *redis.Client
}

var MyRedis redisHandle

func Init(addrs []string, passwd string) error {
    if 1 == len(addrs) {
        MyRedis.IsCluster = false
        MyRedis.OneNodeClient = redis.NewClient(&redis.Options{
            Addr: addrs[0],
            Password: passwd,
            DB: 0,
            MaxRetries: 2,
            DialTimeout: time.Second,
            ReadTimeout: time.Second,
            WriteTimeout: time.Second,
            PoolSize: 8,
            MinIdleConns: 4,
            MaxConnAge: time.Second*30,
            PoolTimeout: time.Second,
            IdleTimeout: time.Second*5,
            IdleCheckFrequency: time.Second*10,
        })
        pong, err := MyRedis.OneNodeClient.Ping().Result()
    } else {
        MyRedis.IsCluster = true
        MyRedis.ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
            Addrs: addrs,
            Password: passwd,
            MaxRetries: 2,
            DialTimeout: time.Second,
            ReadTimeout: time.Second,
            WriteTimeout: time.Second,
            PoolSize: 8,
            MinIdleConns: 4,
            MaxConnAge: time.Second*30,
            PoolTimeout: time.Second,
            IdleTimeout: time.Second*5,
            IdleCheckFrequency: time.Second*10,
        })
        pong, err := MyRedis.ClusterClient.Ping().Result()
    }
    if err != nil {
        return err
    }
    return nil
}

