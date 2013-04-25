package redis

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	Remote string
	Psw    string
	Db     int

	conn net.Conn
}

func openConn(remote, psw string, db int) (net.Conn, error) {
	conn, err := net.Dial("tcp", remote)
	if err != nil {
		return nil, err
	}

	bufrd := bufio.NewReader(conn)
	if psw != "" {
		// if the password was given, do authentication
		_, err = conn.Write([]byte("AUTH " + psw + "\r\n"))
		if err != nil {
			return nil, err
		}
		_, err = readResponse(bufrd)
	}

	if db != 0 {
		// if the database number was given, do selection
		_, err = conn.Write([]byte("SELECT " + strconv.Itoa(db) + "\r\n"))
		if err != nil {
			return nil, err
		}
		_, err = readResponse(bufrd)
	}
	return conn, nil
}

func readLine(bufrd *bufio.Reader) ([]byte, error) {
	p, err := bufrd.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("REDISGO: long response line")
	}
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, errors.New("REDISGO: bad response line terminator")
	}
	return p[:i], nil
}

func readResponse(bufrd *bufio.Reader) (interface{}, error) {
	line, err := readLine(bufrd)
	// fmt.Println(string(line)) // debug
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("redigo: short response line")
	}

	// command executed successfully, return "+OK"
	if line[0] == '+' {
		ret := line[1:]
		return ret, nil
	}

	// command executed failed, return "-ERR ..."
	if bytes.HasPrefix(line, []byte("-ERR")) {
		errmsg := line[5:]
		return nil, errors.New("REDISGO: " + string(errmsg))
	}

	// followed by a number
	if line[0] == ':' {
		num, err := strconv.Atoi(string(line[1:]))
		return num, err
	}

	if line[0] == '$' {
		n, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		if n < 0 {
			return make([]byte, 0), nil
		}

		p := make([]byte, n)
		if _, err = io.ReadFull(bufrd, p); err != nil {
			return nil, err
		}
		line, err := readLine(bufrd)
		if err != nil {
			return nil, err
		}
		if len(line) != 0 {
			return nil, errors.New("REDISGO: bad bulk format")
		}

		ret := unquote(p)
		return ret, nil
	}

	if line[0] == '*' {
		n, err := strconv.Atoi(string(line[1:]))
		if err != nil || n < 0 {
			return nil, err
		}
		r := make([][]byte, n)
		for i := range r {
			rt, err := readResponse(bufrd)
			r[i] = rt.([]byte)
			if err != nil {
				return nil, err
			}
		}
		return r, nil
	}

	err = errors.New("REDISGO: Unkown reply message") // uncatched type
	return make([]byte, 0), err
}

func sendRecv(conn net.Conn, args ...string) (interface{}, error) {
	if conn == nil {
		return nil, errors.New("REDISGO: Connection is not opened yet!")
	}

	c := strings.Join(args, " ")
	c += "\r\n"
	// fmt.Println(c) // debug

	_, err := conn.Write([]byte(c))
	if err != nil {
		return nil, err
	}

	bufrd := bufio.NewReader(conn)
	r, err := readResponse(bufrd)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (client *Client) Connect() error {
	var err error
	if client.conn, err = openConn(client.Remote, client.Psw, client.Db); err != nil {
		return err
	}
	return nil
}

func (client *Client) Disconnect() {
	client.conn = nil
}

func (client *Client) IsActive() bool {
	return (client.conn != nil)
}

// General Commands
func (client *Client) Select(db int) error {
	_, err := sendRecv(client.conn, "SELECT", strconv.Itoa(db))
	return err
}

func (client *Client) Set(key []byte, value []byte) error {
	_, err := sendRecv(client.conn, "SET", string(quote(key)), string(quote(value)))
	return err
}

func (client *Client) Get(key []byte) ([]byte, error) {
	r, err := sendRecv(client.conn, "GET", string(quote(key)))
	if err != nil {
		return nil, err
	}

	return r.([]byte), nil
}

func (client *Client) Keys(arg []byte) ([][]byte, error) {
	r, err := sendRecv(client.conn, "KEYS", string(quote(arg)))
	if err != nil {
		return nil, err
	}

	ret := r.([][]byte)
	return ret, nil
}

func (client *Client) Hmset(key []byte, arg map[string][]byte) error {
	c := make([]string, 0)
	c = append(c, "HMSET", string(quote(key)))
	for k, v := range arg {
		c = append(c, k, string(quote(v)))
	}

	c = append(c, "\r\n")
	_, err := sendRecv(client.conn, c...)
	return err
}

func (client *Client) Hmget(key []byte, field ...[]byte) ([][]byte, error) {
	c := make([]string, 0)
	c = append(c, "HMGET", string(quote(key)))
	for _, f := range field {
		c = append(c, string(quote(f)))
	}

	r, err := sendRecv(client.conn, c...)
	if err != nil {
		return make([][]byte, 0), err
	}
	return r.([][]byte), nil
}

func (client *Client) Sadd(key []byte, members ...[]byte) (int, error) {
	c := make([]string, 0)
	c = append(c, "SADD", string(quote(key)))
	for _, m := range members {
		c = append(c, string(quote(m)))
	}

	num, err := sendRecv(client.conn, c...)
	return num.(int), err
}

func (client *Client) Smembers(key []byte) ([][]byte, error) {
	r, err := sendRecv(client.conn, "SMEMBERS", string(quote(key)))
	if err != nil {
		return nil, err
	}
	return r.([][]byte), nil
}
