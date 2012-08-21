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
    if err := client.Set("foo", "bar"); err != nil {
        fmt.Println("SET" + err.Error())
    }
}

func TestGet() {
    str1 := "bar"
    str2, err := client.Get("foo")

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
    str, err := client.Keys("nil*")
    if err != nil {
        fmt.Println("KEYS" + err.Error())
    }
    fmt.Println("KEYS:")
    fmt.Println(str)
    fmt.Println("----------------")
}

func TestHmset() {
    m := map[string]string{
        "key1" : "value1",
        "key2" : "value2",
    }
    if err := client.Hmset("hmset", m); err != nil {
        fmt.Println("HMSET" + err.Error())
    }
}

func TestHmget() {
    res, err := client.Hmget("hmset", "key9", "key2")
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("HMGET:")
    for _, str := range res {
        fmt.Printf("%s ", str)
    }
    fmt.Println()
    fmt.Println("----------------")
}

func TestSadd() {
    num, err := client.Sadd("sadd", "a", "b", "c", "d", "e")
    if err != nil {
        fmt.Println("SADD" + err.Error())
    }
    fmt.Println("SADD:")
    fmt.Println(num)
    fmt.Println("----------------")
}

func TestSmembers() {
    str, err := client.Smembers("sadd")
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
