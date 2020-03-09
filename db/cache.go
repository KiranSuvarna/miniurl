package db

import (
	"bytes"
	"crypto/tls"
	"strconv"
	"time"

	"bitbucket.org/smartclean/routines-go/config"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

// RedisConn Creates a new Redis client
type RedisConn struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
	expiration    time.Duration
}

// NewRedis Creates a new Redis connection
func NewRedis(conf *config.RedisConfig) (*RedisConn, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       conf.Address,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &RedisConn{
		client:        client,
		clusterClient: nil,
		expiration:    conf.Expiration,
	}, nil
}

//NewRedisCluster Creates a new Redis connection to cluster
func NewRedisCluster(conf *config.RedisClusterConfig) (*RedisConn, error) {
	cClient := redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots: func() ([]redis.ClusterSlot, error) {
			slots := []redis.ClusterSlot{
				{
					Start: 0,
					End:   8191,
					Nodes: []redis.ClusterNode{{
						Addr: conf.Master,
					}, {
						Addr: conf.Replica,
					}},
				},
			}
			return slots, nil
		},
		RouteRandomly: true,
		Password:      conf.Password,
		TLSConfig:     &tls.Config{},
		MaxRetries:    conf.MaxRetries,
	})

	return &RedisConn{
		clusterClient: cClient,
		expiration:    conf.Expiration,
	}, nil
}

// Close Closes the Redis connection
func (c *RedisConn) Close() error {
	if c.client != nil {
		return c.client.Close()
	}

	return nil
}

// Set Sets a value in cache
func (c *RedisConn) Set(key string, value interface{}) error {
	return c.SetWithExpiration(key, value, c.expiration)
}

// SetWithExpiration Sets a value in cache with expiration
func (c *RedisConn) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	w := new(bytes.Buffer)

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(w).Encode(value); err != nil {
		return err
	}

	// log.WithFields(log.Fields{"key": key, "input": value, "output": w.String()}).Info("cached")
	return c.client.Set(key, w.String(), expiration).Err()
}

// Incr Increments value
func (c *RedisConn) Incr(key string) error {
	_, err := c.client.Incr(key).Result()

	return err
}

// Get Gets value from cache
func (c *RedisConn) Get(key string) (string, error) {
	res, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return res, nil
}

// GetMultiple Gets multiple values from cache
func (c *RedisConn) GetMultiple(keys []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, k := range keys {
		results[k] = ""
	}

	res, err := c.client.MGet(keys...).Result()
	if err == redis.Nil {
		return results, nil
	} else if err != nil {
		return results, err
	}

	for i, k := range keys {
		if res[i] != nil {
			results[k] = res[i].(string)
		}
	}

	return results, nil
}

// GetInt Gets integer value from cache
func (c *RedisConn) GetInt(key string) (int64, error) {
	res, err := c.Get(key)
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return strconv.ParseInt(res, 10, 64)
}

// GetIntMultiple Gets multiple integer values from cache
func (c *RedisConn) GetIntMultiple(keys []string) (map[string]int64, error) {
	results := make(map[string]int64, 0)
	for _, k := range keys {
		results[k] = 0
	}

	res, err := c.client.MGet(keys...).Result()
	if err == redis.Nil {
		return results, nil
	} else if err != nil {
		return results, err
	}

	for i, k := range keys {
		if res[i] != nil {
			s := res[i].(string)
			if s != "" {
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					results[k] = v
				}
			}
		}
	}

	return results, nil
}

func (c *RedisConn) Del(key string) (int64, error) {
	res, err := c.client.Del(key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return res, nil
}

// Keys Gets all the matching keys from cache
func (c *RedisConn) Keys(pattern string) ([]string, error) {
	res, err := c.client.Keys(pattern).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return res, nil
}