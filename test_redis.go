package main

import (
    "fmt"
    "redisgo"
)

var cli redisgo.Client

func cliInit() {
    cli.Remote = "127.0.0.1:6379"
    // cli.Psw = ""
    cli.Db = 13

    if err := cli.Connect(); err != nil {
        fmt.Println("Init: " + err.Error())
    }
}

func TestSet() {
    if err := cli.Set([]byte("foo"), []byte("bar")); err != nil {
        fmt.Println("SET" + err.Error())
    }
}

func TestGet() {
    str1 := "bar"
    byte2, err := cli.Get([]byte("foo"))
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
    str, err := cli.Keys([]byte("*"))
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
    if err := cli.Hmset([]byte("hmset"), m); err != nil {
        fmt.Println("HMSET" + err.Error())
    }
}

func TestHmget() {
    res, err := cli.Hmget([]byte("hmset"), []byte("key1"), []byte("key2"))
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
    num, err := cli.Sadd([]byte("sadd"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"))
    if err != nil {
        fmt.Println("SADD" + err.Error())
    }
    fmt.Println("SADD:")
    fmt.Println(num)
    fmt.Println("----------------")
}

func TestSmembers() {
    str, err := cli.Smembers([]byte("sadd"))
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("SMEMBERS:")
    fmt.Println(str)
    fmt.Println("----------------")
}

func main() {
    cliInit()
    defer cli.Disconnect()

    TestSet()
    TestGet()
    TestKeys()
    TestHmset()
    TestHmget()
    TestSadd() 
    TestSmembers()
}
