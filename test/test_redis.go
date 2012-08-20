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
        fmt.Println(err)
    }
}

func TestSet() {
    if err := client.Set("foo", "bar"); err != nil {
        fmt.Println(err)
    }
}

func TestGet() {
    str1 := "bar"
    str2, err := client.Get("foo")

    if err != nil {
        fmt.Println(err)
    }
    if str2 != nil && str2 != nil && str1 != *str2 {
        fmt.Println("str1:", str1)
        fmt.Println("str2:", *str2)
        fmt.Println("str unequal")
    }
}

func TestKeys() {
    str, err := client.Keys("*")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(str)
}

func TestHmset() {
    m := map[string]string{
        "key1" : "value1",
        "key2" : "value2",
    }
    if err := client.Hmset("hmset", &m); err != nil {
        fmt.Println(err)
    }
}

func TestHmget() {
    res, err := client.Hmget("hmset", "key1", "key2")
    if err != nil {
        fmt.Println(err)
    }
    for _, pstr := range res {
        if pstr != nil {
            fmt.Printf("%s ", *pstr)
        } else {
            fmt.Printf("(nil) ")
        }
    }
    fmt.Println()
}

func TestSadd() {
    num, err := client.Sadd("sadd", "a", "b", "c", "d", "e")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(num)
}

func TestSmembers() {
    str, err := client.Smembers("sadd")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(str)
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
