package main

import (
	"fmt"
	"strconv"
	"time"

	sw "week05/slidingwindow"

	"github.com/go-redis/redis"
)

func main() {
	size := time.Second
	store := NewRedisDatastore(
		redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		2*size,
	)

	lim, stop := sw.NewLimiter(size, 1, func() (sw.Window, sw.StopFunc) {
		return sw.NewSyncWindow("test", sw.NewNonblockingSynchronizer(store, 500*time.Millisecond))
	})
	defer stop()

	ok := lim.Allow()
	fmt.Printf("ok: %v\n", ok)
	ok = lim.Allow()
	fmt.Printf("ok: %v\n", ok)
}

type RedisDatastore struct {
	client redis.Cmdable
	ttl    time.Duration
}

func NewRedisDatastore(client redis.Cmdable, ttl time.Duration) *RedisDatastore {
	return &RedisDatastore{client: client, ttl: ttl}
}

func (d *RedisDatastore) fullKey(key string, start int64) string {
	return fmt.Sprintf("%s@%d", key, start)
}

func (d *RedisDatastore) Add(key string, start, value int64) (int64, error) {
	k := d.fullKey(key, start)
	c, err := d.client.IncrBy(k, value).Result()
	if err != nil {
		return 0, err
	}
	d.client.Expire(k, d.ttl).Result()
	return c, err
}

func (d *RedisDatastore) Get(key string, start int64) (int64, error) {
	k := d.fullKey(key, start)
	value, err := d.client.Get(k).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}
