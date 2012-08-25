package redis_test

import (
	"."
	"fmt"
	"runtime"
	"testing"
)

var cli redis.Client

func init() {
	runtime.GOMAXPROCS(1)

	cli.Remote = "127.0.0.1:6379"
	// cli.Psw = "lucky" // optional
	cli.Db = 13 // optional

	if err := cli.Connect(); err != nil {
		panic(err.Error())
	}
}

func TestSet(t *testing.T) {
	in := []byte{0x5A, 0x20, 0x5A}
	if err := cli.Set([]byte("foo"), in); err != nil {
		t.Fatal("SET failed", err)
	}
}

func ExampleGet() {
	ret, err := cli.Get([]byte("foo"))
	if err != nil {
		panic("GET failed")
	}
	fmt.Println(string(ret))
	// Output: Z Z
}

func TestKeys(t *testing.T) {
	_, err := cli.Keys([]byte("*"))
	if err != nil {
		t.Fatal("KEYS * failed", err)
	}
}

func TestHmset(t *testing.T) {
	m := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	if err := cli.Hmset([]byte("hmset"), m); err != nil {
		t.Fatal("HMSET failed", err)
	}
}

func TestHmget(t *testing.T) {
	ret, err := cli.Hmget([]byte("hmset"), []byte("key1"), []byte("key2"))
	if err != nil {
		t.Fatal("HMGET failed", err)
	}
	if string(ret[0]) != "value1" || string(ret[1]) != "value2" {
		t.Fatal("HMGET failed", err)
	}
}

func TestSadd(t *testing.T) {
	_, err := cli.Sadd([]byte("sadd"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"))
	if err != nil {
		t.Fatal("SADD failed", err)
	}
}

func TestSmembers(t *testing.T) {
	_, err := cli.Smembers([]byte("smembers"))
	if err != nil {
		t.Fatal("SMEMBERS failed", err)
	}
}
