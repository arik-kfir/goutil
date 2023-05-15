package redisutil

import (
	"crypto/tls"
	"fmt"
	"github.com/rueian/rueidis"
	"strconv"
)

type RedisConfig struct {
	Host     string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" required:"yes"`
	Port     int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	TLS      bool   `env:"TLS" value-name:"TLS" long:"tls" description:"Whether to use TLS to connect to Redis"`
	PoolSize int    `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
}

func (c *RedisConfig) Connect(clientName string) (rueidis.Client, error) {
	redisClientOption := rueidis.ClientOption{
		InitAddress:      []string{c.Host + ":" + strconv.Itoa(c.Port)},
		ClientName:       clientName,
		BlockingPoolSize: c.PoolSize,
	}
	if c.TLS {
		redisClientOption.TLSConfig = &tls.Config{ServerName: c.Host}
	}
	redisClient, err := rueidis.NewClient(redisClientOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return redisClient, nil
}
