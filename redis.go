package redisgo

import (
    // "fmt"
    "net"
    "strconv"
    "bytes"
    "encoding/binary"
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

type RError string

func (err RError) Error() string {
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

func readResponse(conn net.Conn) (interface{}, error) {
    var data []byte = make([]byte, bufSize)

    n, err := conn.Read(data)
    if err != nil {
        return make([]byte, 0), err
    }

    line := bytes.TrimSpace(data[0:n])

    if line[0] == '+' {
        res := line[1:]
        return res, nil
    }

    if bytes.HasPrefix(line, []byte("-ERR")) {
        errmsg := line[5:]
        return nil, RError(errmsg)
    }

    if line[0] == ':' {
        num, err := strconv.Atoi(string(line[1:]))
        return num, err
    }

    if line[0] == '$' {
        if bytes.HasPrefix(line, []byte("$-1")) {
            return make([]byte, 0), nil
        }

        list := bytes.Split(line, []byte("\r\n")) // fmt.Printf("list: %v", list) // Debug
        res := list[1]
        return res, nil
    }

    if line[0] == '*' {
        list := bytes.Split(line, []byte("\r\n"))
        // fmt.Printf("list: %v\n", list) // debug

        nsize, err := strconv.Atoi(string(list[0][1:]))
        if err != nil {
            return make([]byte, 0), err
        }

        rbyte := make([][]byte, nsize)
        var k int = 0
        for i := 1; i < len(list); i++ {
            if (bytes.HasPrefix(list[i], []byte("$-1"))) {
                k += 1;
                continue
            }

            i += 1
            rbyte[k] = list[i]
            k += 1
        }
        return rbyte, nil
    }

    err = RError("Unkown reply message") // uncatched type
    return make([]byte, 0), err
}

func sendRecv(conn net.Conn, args ...[]byte) (interface{}, error) {
    if conn == nil {
        return nil, RError("connection is not created yet!")
    }

    c := bytes.Join(args, []byte(" "))
    c = append(c, []byte("\r\n")...)
    _, err := conn.Write(c)
    if err != nil {
        return nil, err
    }

    r, err := readResponse(conn)
    if err != nil {
        return nil, err
    }
    return r, nil
}

func (client *Client) Connect() error {
    var err error
    client.conn, err = openConn(client.Remote, client.Psw, client.Db)
    return err
}

func (client *Client) Disconnect() {
    client.conn = nil
}

// General Commands
func (client *Client) Select(db int) error {
    bin := make([]byte, 4)
    binary.BigEndian.PutUint32(bin, uint32(db))
    _, err := sendRecv(client.conn, []byte("SELECT"), bin)
    return err
}

func (client *Client) Set(key []byte, value []byte) error {
    _, err := sendRecv(client.conn, []byte("SET"), key, value)
    return err
}

func (client *Client) Get(key []byte) ([]byte, error) {
    r, err := sendRecv(client.conn, []byte("GET"), key)
    if err != nil {
        return nil, err
    }

    return r.([]byte), nil
}

func (client *Client) Keys(arg []byte) ([][]byte, error) {
    r, err := sendRecv(client.conn, []byte("KEYS"), arg)
    if err != nil {
        return nil, err
    }

    if r == nil {
        return make([][]byte, 0), nil
    }
    res := r.([][]byte)
    return res, nil
}

func (client *Client) Hmset(key []byte, arg map[string][]byte) (error) {
    c := make([][]byte, 0)
    c = append(c, []byte("HMSET"), key)
    for k, v := range arg {
        c = append(c, []byte(k), v)
    }

    c = append(c, []byte("\r\n"))
    _, err := sendRecv(client.conn, c...)
    return err
}

func (client *Client) Hmget(key []byte, fields ...[]byte) ([][]byte, error) {
    c := make([][]byte, 0)
    c = append(c, []byte("HMGET"), key)
    c = append(c, fields...)

    r, err := sendRecv(client.conn, c...)
    if err != nil {
        return nil, err
    }
    return r.([][]byte), nil
}

func (client *Client) Sadd(key []byte, members ...[]byte) (int, error) {
    c := make([][]byte, 0)
    c = append(c, []byte("SADD"), key)
    c = append(c, members...)

    num, err := sendRecv(client.conn, c...)
    if err != nil {
        return 0, err
    }
    return num.(int), err
}

func (client *Client) Smembers(key []byte) ([][]byte, error) {
    r, err := sendRecv(client.conn, []byte("SMEMBERS"), key)
    if err != nil {
        return nil, err
    }
    return r.([][]byte), nil
}
