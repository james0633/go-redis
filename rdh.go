package rdh

import (
    "reflect"
    "time"
    "github.com/go-redis/redis"
)

type redisHandle struct {
    IsCluster       bool
    ClusterClient   *redis.ClusterClient
    OneNodeClient   *redis.Client
}

var MyRedis redisHandle
var RcNames map[string]reflect.Value

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
        vf := reflect.ValueOf(MyRedis.OneNodeClient)
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
        vf := reflect.ValueOf(MyRedis.ClusterClient)
        pong, err := MyRedis.ClusterClient.Ping().Result()
    }
    if err != nil && pong != "PONG" {
        return err
    }
    var vm string
    for i := 0; i < vf.NumMethod(); i++ {
		vm = vf.Method(i).Name
		RcNames[vm] = vf.Method(i)
	}
    return nil
}

func Call(name string, params ... interface{}) (result []reflect.Value, err error) {
    if _, ok := RcNames[name]; !ok {
        return nil, errors.New(name + " does not exist.")
    }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    result = RcNames[name].Call(in)
    return
}
