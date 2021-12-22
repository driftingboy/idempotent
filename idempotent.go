package idempotent

import (
	"errors"
	"fmt"

	redis "github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

type Idempotence struct {
	storage Storage
}

// generate a idempotence string
func GenerateID() string {
	return uuid.New().String()
}

// If you need to use other storage methods, or want to use the existing client,
// you can implement 'Storage' interface，use function 'New(s Storage) *idempotence'.
func New(s Storage) *Idempotence {
	return &Idempotence{
		storage: s,
	}
}

// config is not nil
func NewWithOpts(config *Config, opts ...Option) *Idempotence {

	redisOpts := &redis.ClusterOptions{
		Addrs:    config.RedisAddrs,
		Password: config.Password,
	}

	for _, o := range opts {
		o(redisOpts)
	}

	cc := redis.NewClusterClient(redisOpts)
	fmt.Printf("ping %+v", cc.Do("ping").String())
	return &Idempotence{
		storage: newStorage(cc),
	}
}

// exist return false, if no exist, set and return true.
// concurrent safe
func (i *Idempotence) CheckIdempotence(ID string) (bool, error) {
	if i.storage == nil {
		return false, errors.New("redisCluster can not be nil!")
	}
	// 考虑设置 value=time
	return i.storage.SetNX(ID, true)
}

// delete idempotence ID, if err != nil, means equipment failure
func (i *Idempotence) Delete(ID string) error {
	return i.storage.Delete(ID)
}
