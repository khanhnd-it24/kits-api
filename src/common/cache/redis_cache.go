package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kits/api/src/common/configs"
)

func NewRedisClient(cf *configs.Config) (redis.UniversalClient, error) {
	hosts := cf.Redis.Hosts
	if len(hosts) == 0 {
		return nil, errors.New("no such host redis")
	}
	var client redis.UniversalClient
	isClusterMode := len(hosts) > 1
	if isClusterMode {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hosts,
			Username: cf.Redis.Username,
			Password: cf.Redis.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     hosts[0],
			Username: cf.Redis.Username,
			Password: cf.Redis.Password,
		})
	}
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis %w", err)
	}

	return client, nil
}
