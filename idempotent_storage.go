package idempotent

import redis "github.com/go-redis/redis/v7"

type Storage interface {
	// if exist return false
	// if not exist, set value and return true
	// need concurrent-safe
	SetNX(id string, value interface{}) (bool, error)
	// delete idempotence flag
	Delete(id string) error
}

// default impl
type redisClusterStorage struct {
	redisCluster *redis.ClusterClient
}

// if you had a ClusterClient, use function ` New()`
// 此处依赖了其他库，如果外部业务已有redis client，但不是 go-redis 这个库就无法使用
func newStorage(rc *redis.ClusterClient) Storage {
	return &redisClusterStorage{
		redisCluster: rc,
	}
}

func (r *redisClusterStorage) SetNX(id string, value interface{}) (bool, error) {
	return r.redisCluster.SetNX(id, value, 0).Result()
}

func (r *redisClusterStorage) Delete(id string) error {
	return r.redisCluster.Del(id).Err()
}
