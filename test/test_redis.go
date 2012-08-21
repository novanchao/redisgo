package main

import (
    "fmt"
    "goredis"
)

var client redis.Client

func clientInit() {
    client.Remote = "202.119.236.131:6379"
    client.Psw = "redis"
    client.Db = 13

    if err := client.Connect(); err != nil {
        fmt.Println("Init: " + err.Error())
    }
}

func TestSet() {
    if err := client.Set([]byte("foo"), []byte("bar")); err != nil {
        fmt.Println("SET" + err.Error())
    }
}

func TestGet() {
    str1 := "bar"
    byte2, err := client.Get([]byte("foo"))
    str2 := string(byte2)

    if err != nil {
        fmt.Println("GET" + err.Error())
    }
    if str1 != str2 {
        fmt.Println("GET:")
        fmt.Println("str1:", str1)
        fmt.Println("str2:", str2)
        fmt.Println("str unequal")
        fmt.Println("----------------")
    }
}

func TestKeys() {
    str, err := client.Keys([]byte("*"))
    if err != nil {
        fmt.Println("KEYS" + err.Error())
    }
    fmt.Println("KEYS:")
    fmt.Println(str)
    fmt.Println("----------------")
}

func TestHmset() {
    m := map[string][]byte{
        "key1" : []byte("value1"),
        "key2" : []byte("value2"),
    }
    if err := client.Hmset([]byte("hmset"), m); err != nil {
        fmt.Println("HMSET" + err.Error())
    }
}

func TestHmget() {
    res, err := client.Hmget([]byte("hmset"), []byte("key9"), []byte("key2"))
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("HMGET:")
    for _, str := range res {
        fmt.Printf("%s ", string(str))
    }
    fmt.Println()
    fmt.Println("----------------")
}

func TestSadd() {
    num, err := client.Sadd([]byte("sadd"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"))
    if err != nil {
        fmt.Println("SADD" + err.Error())
    }
    fmt.Println("SADD:")
    fmt.Println(num)
    fmt.Println("----------------")
}

func TestSmembers() {
    str, err := client.Smembers([]byte("sadd"))
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("SMEMBERS:")
    fmt.Println(str)
    fmt.Println("----------------")
}

func main() {
    clientInit()
    defer client.Disconnect()

    TestSet()
    TestGet()
    TestKeys()
    TestHmset()
    TestHmget()
    TestSadd() 
    TestSmembers()
}
