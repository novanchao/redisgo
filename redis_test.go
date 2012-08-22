package redisgo_test

import (
	"."
	"testing"
)

var cli redisgo.Client

func TestSet(t *testing.T) {
	if err := cli.Set([]byte("foo"), []byte("bar")); err != nil {
		t.Log("SET" + err.Error())
	}
}

func TestGet(t *testing.T) {
	str1 := "bar"
	byte2, err := cli.Get([]byte("foo"))
	str2 := string(byte2)

	if err != nil {
		t.Log("GET" + err.Error())
	}
	if str1 != str2 {
		t.Log("GET:")
		t.Log("str1:", str1)
		t.Log("str2:", str2)
		t.Log("str unequal")
		t.Log("----------------")
	}
}

func TestKeys(t *testing.T) {
	str, err := cli.Keys([]byte("*"))
	if err != nil {
		t.Log("KEYS" + err.Error())
	}
	t.Log("KEYS:")
	t.Log(str)
	t.Log("----------------")
}

func TestHmset(t *testing.T) {
	m := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}
	if err := cli.Hmset([]byte("hmset"), m); err != nil {
		t.Log("HMSET" + err.Error())
	}
}

func TestHmget(t *testing.T) {
	res, err := cli.Hmget([]byte("hmset"), []byte("key1"), []byte("key2"))
	if err != nil {
		t.Log(err.Error())
	}
	t.Log("HMGET:")
	for _, str := range res {
		t.Logf("%s ", string(str))
	}
	t.Log()
	t.Log("----------------")
}

func TestSadd(t *testing.T) {
	num, err := cli.Sadd([]byte("sadd"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"))
	if err != nil {
		t.Log("SADD" + err.Error())
	}
	t.Log("SADD:")
	t.Log(num)
	t.Log("----------------")
}

func TestSmembers(t *testing.T) {
	str, err := cli.Smembers([]byte("sadd"))
	if err != nil {
		t.Log(err.Error())
	}
	t.Log("SMEMBERS:")
	t.Log(str)
	t.Log("----------------")
}

func init() {
	cli.Remote = "127.0.0.1:6379"
	// cli.Psw = ""
	cli.Db = 13

	if err := cli.Connect(); err != nil {
		println("Init: " + err.Error())
	} else {
		defer cli.Disconnect()
	}
}
