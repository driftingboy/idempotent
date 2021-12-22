package idempotent_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/driftingboy/idempotent"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

func Test_Use_DefultImpl(t *testing.T) {
	idem := idempotent.NewWithOpts(&idempotent.Config{
		RedisAddrs: []string{"127.0.0.1:6379"},
		Password:   "142589",
	})

	id := idempotent.GenerateID()

	// 使用方式一
	for i := 0; i < 2; i++ {
		if ok, err := idem.CheckIdempotence(id); err != nil {
			t.Fatal(err)
			return
		} else if !ok {
			t.Log("idempotent exist")
			return
		}
	}

	// do bussiness
	t.Log("exec bussiness")

	// 使用方式二
	// bussFunc := func() error {
	// 	// ...
	// }
	// i.DoWithIdempotent(bussFunc)

}

// 测试自定义 storage 逻辑

type RedisStorage struct {
	client redis.Client
}

func (r *RedisStorage) SetNX(id string, value interface{}) (bool, error) {
	return r.client.SetNX(id, value, 0).Result()
}

func (r *RedisStorage) Delete(id string) error {
	_, err := r.client.Del(id).Result()
	return err
}

var redisStorage = &RedisStorage{
	client: *redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "142589",
	}),
}

func Example_Custom_Storage(t *testing.T) {

	idem := idempotent.New(redisStorage)

	id1 := idempotent.GenerateID()

	for i := 0; i < 2; i++ {
		err := MonitorUseFunc(idem, id1)
		assert.NoError(t, err)
	}

	err := idem.Delete(id1)
	assert.NoError(t, err)
	fmt.Println("had delete idempotence id.")

	err = MonitorUseFunc(idem, id1)
	assert.NoError(t, err)
	//Output:
	//do bussiness ...
	//had the idempotence id, jump over bussiness.
	//had delete idempotence id.
	//do bussiness ...
}

// 并发测试

func Test_Concurrency(t *testing.T) {
	idem := idempotent.New(redisStorage)

	var (
		id1          = idempotent.GenerateID()
		existCount   = 0
		successCount = 0
	)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ok, err := idem.CheckIdempotence(id1); err != nil {
				return
			} else if !ok {
				existCount++
				return
			}
			successCount++
		}()
	}

	wg.Wait()

	err := idem.Delete(id1)
	assert.NoError(t, err)

	assert.Equal(t, 1, successCount)
	assert.Equal(t, 9, existCount)

}

func MonitorUseFunc(idem *idempotent.Idempotence, id string) error {
	if ok, err := idem.CheckIdempotence(id); err != nil {
		return err
	} else if !ok {
		fmt.Println("had the idempotence id, jump over bussiness.")
		return nil
	}

	fmt.Println("do bussiness ...")
	return nil
}
