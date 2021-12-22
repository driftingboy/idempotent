# idempotent
go idempotent tool

# demo
```go
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

}
```

详细使用参考 idempotent_test.go
