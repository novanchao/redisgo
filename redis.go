package redis

import (
    // "fmt"
    "net"
    "strconv"
    "strings"
)

const (
    bufSize = 1024
)

type Client struct {
    Remote string
    Psw string
    Db int
    conn net.Conn
}

type RedisError string

func (err RedisError) Error() string {
    return "REDIS ERROR: " + string(err)
}

func openConn(remote string, psw string, db int) (net.Conn, error) {
    conn, err := net.Dial("tcp", remote)
    if err != nil {
        return nil, err
    }

    if psw != "" {
        // if the password was given, do authentication
        _, err = conn.Write([]byte("AUTH " + psw + "\r\n"))
        if err != nil {
            return nil, err
        }
        _, err = readResponse(conn)
    }

    if db!= 0 {
        _, err = conn.Write([]byte("SELECT " + strconv.Itoa(db) + "\r\n"))
        if err != nil {
            return nil, err
        }
        _, err = readResponse(conn)
    }
    return conn, err
}

func sendRecv(conn net.Conn, args ...string) (interface{}, error) {
    cmd := strings.Join(args, " ")
    if conn == nil {
        return nil, RedisError("connection is not created yet!")
    }

    _, err := conn.Write([]byte(cmd + "\r\n"))
    if err != nil {
        return nil, err
    }

    r, err := readResponse(conn)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func readResponse(conn net.Conn) (interface{}, error) {
    var data []byte = make([]byte, bufSize)

    n, err := conn.Read(data)
    if err != nil {
        return nil, err
    }

    line := strings.TrimSpace(string(data[0:n]))
    // fmt.Println("line:", line) // Debug
    // fmt.Println("len(line):", len(line)) // Debug

    if line[0] == '+' {
        res := line[1:]
        return &res, nil
    }

    if strings.HasPrefix(line, "-ERR ") {
        errmsg := line[5:]
        return nil, RedisError(errmsg)
    }

    if line[0] == ':' {
        num, err := strconv.Atoi(line[1:len(line)])
        return num, err
    }

    if line[0] == '$' {
        if strings.HasPrefix(line, "$-1") {
            return nil, nil
        }

        list := strings.Split(line, "\r\n") // fmt.Printf("list: %v", list) // Debug
        res := list[1]
        return &res, nil
    }

    if line[0] == '*' {
        list := strings.Split(line, "\r\n")
        // fmt.Printf("list: %v\n", list) // debug

        nsize, err := strconv.Atoi(list[0][1:])
        if err != nil {
            return nil, err
        }
        var k int = 0
        reslice := make([]*string, nsize)
        // fmt.Println(strconv.Itoa(nsize), strconv.Itoa(len(list))) // debug
        for i := 1; i < len(list); i++ {
            if (strings.HasPrefix(list[i], "$-1")) {
                // TODO: how to deal with the string who has "$-1"
                // TODO: need error or not
                continue
            }

            i += 1
            reslice[k] = &list[i]
            k += 1
        }
        return reslice, nil
    }

    err = RedisError("Unkown reply message") // uncatched type
    return nil, err
}

func (client *Client) Connect() error {
    var err error
    client.conn, err = openConn(client.Remote, client.Psw, client.Db)
    return err
}

func (client *Client) Disconnect() error {
    var err error
    client.conn, err = openConn(client.Remote, client.Psw, client.Db)
    return err
}

// General Commands
func (client *Client) Select(db int) error {
    _, err := sendRecv(client.conn, "SELECT", strconv.Itoa(db))
    return err
}

func (client *Client) Set(key string, value string) error {
    _, err := sendRecv(client.conn, "SET", key, value)
    return err
}

func (client *Client) Get(key string) (*string, error) {
    r, err := sendRecv(client.conn, "GET", key)
    if err != nil {
        return nil, err
    }

    var res *string
    if r != nil {
        res = r.(*string)
    }
    return res, nil
}

func (client *Client) Keys(arg string) ([]string, error) {
    r, err := sendRecv(client.conn, "KEYS", arg)
    if err != nil {
        return nil, err
    }

    res := r.([]*string)
    strs := make([]string, len(res))
    for i, _ := range res {
        strs[i] = *res[i]
    }
    return strs, nil
}

func (client *Client) Hmset(key string, arg *map[string]string) (error) {
    cmd := "HMSET " + key
    for k, v := range *arg {
        cmd += " " + k + " " + v
    }

    cmd += "\r\n"
    _, err := sendRecv(client.conn, cmd)
    return err
}

func (client *Client) Hmget(key string, fields ...string) ([]*string, error) {
    s := []string{"HMGET", key}
    cmd := append(s, fields...)

    r, err := sendRecv(client.conn, cmd...)
    if err != nil {
        return nil, err
    }
    return r.([]*string), nil
}

func (client *Client) Sadd(key string, members ...string) (int, error) {
    s := []string{"SADD", key}
    cmd := append(s, members...)

    num, err := sendRecv(client.conn, cmd...)
    if err != nil {
        return 0, err
    }
    return num.(int), err 
}

func (client *Client) Smembers(key string) ([]string, error) {
    r, err := sendRecv(client.conn, "SMEMBERS", key)
    if err != nil {
        return nil, err
    }

    res := r.([]*string)
    strs := make([]string, len(res))
    for i, _ := range res {
        strs[i] = *res[i]
    }
    return strs, nil

}
