package redis_test

import (
	"."
	"fmt"
	"runtime"
	"testing"
)

var pool redis.RedisPool

func init() {
	runtime.GOMAXPROCS(5)

	pool.Remote = "127.0.0.1:6379"
	// cli.Psw = "lucky" // optional
	pool.Db = 13 // optional

	pool.CreatePool()
}

func TestPool(t *testing.T) {
	cli, err := pool.PopClient()
	defer pool.PushClient(cli)

	if err != nil {
		t.Fatal(err)
	}

	Set(t, cli)
	Get(cli)
	Keys(t, cli)
	Hmset(t, cli)
	Hmget(t, cli)
	Hmget(t, cli)
	Sadd(t, cli)
	Smembers(t, cli)
}

func Set(t *testing.T, cli *redis.Client) {
	if err := cli.Set([]byte("pool"), []byte("b a\\ \\r")); err != nil {
		t.Fatal("SET failed", err)
	}
}

func Get(cli *redis.Client) {
	ret, err := cli.Get([]byte("pool"))
	if err != nil {
		panic("GET failed")
	}
	fmt.Println(string(ret))
	// Output: b a\ \r
}

func Keys(t *testing.T, cli *redis.Client) {
	_, err := cli.Keys([]byte("*"))
	if err != nil {
		t.Fatal("KEYS * failed", err)
	}
}

func Hmset(t *testing.T, cli *redis.Client) {
	m := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	if err := cli.Hmset([]byte("hmset"), m); err != nil {
		t.Fatal("HMSET failed", err)
	}
}

func Hmget(t *testing.T, cli *redis.Client) {
	ret, err := cli.Hmget([]byte("hmset"), []byte("key1"), []byte("key2"))
	if err != nil {
		t.Fatal("HMGET failed", err)
	}
	if string(ret[0]) != "value1" || string(ret[1]) != "value2" {
		t.Fatal("HMGET failed", err)
	}
}

func Sadd(t *testing.T, cli *redis.Client) {
	_, err := cli.Sadd([]byte("sadd"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"))
	if err != nil {
		t.Fatal("SADD failed", err)
	}
}

func Smembers(t *testing.T, cli *redis.Client) {
	_, err := cli.Smembers([]byte("smembers"))
	if err != nil {
		t.Fatal("SMEMBERS failed", err)
	}
}
