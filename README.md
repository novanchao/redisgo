REDISGO

usage:

    import "redis"

    var cli redis.Clinet

    cli.Remote = "127.0.0.1:6379"
    cli.Psw = "lucky" // optional
    cli.Db = 13 // optional

    if err := cli.Connect(); err != nil {
        // error handle
    }

Good Luck! :)
