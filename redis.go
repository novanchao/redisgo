package redis

import (
	"bytes"
	"net"
	"strconv"
	"strings"
)

const (
	bufSize = 1024
)

type Client struct {
	Remote string
	Psw    string
	Db     int
	conn   net.Conn
}

type RError string

func (err RError) Error() string {
	return "REDIS ERROR: " + string(err)
}

func openConn(remote, psw string, db int) (net.Conn, error) {
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

	if db != 0 {
		// if the database number was given, do selection
		_, err = conn.Write([]byte("SELECT " + strconv.Itoa(db) + "\r\n"))
		if err != nil {
			return nil, err
		}
		_, err = readResponse(conn)
	}
	return conn, nil
}

func readResponse(conn net.Conn) (interface{}, error) {
	var data []byte = make([]byte, bufSize)

	n, err := conn.Read(data)
	if err != nil {
		return make([]byte, 0), err
	}

	line := bytes.TrimSpace(data[0:n])
	// fmt.Println(string(line)) // debug

	// command executed successfully, return "+OK"
	if line[0] == '+' {
		ret := line[1:]
		return ret, nil
	}

	// command executed failed, return "-ERR ..."
	if bytes.HasPrefix(line, []byte("-ERR")) {
		errmsg := line[5:]
		return nil, RError(errmsg)
	}

	// followed by a number
	if line[0] == ':' {
		num, err := strconv.Atoi(string(line[1:]))
		return num, err
	}

	if line[0] == '$' {
		if bytes.HasPrefix(line, []byte("$-1")) {
			return make([]byte, 0), nil
		}

		list := bytes.Split(line, []byte("\r\n"))
		ret := unquote(list[1])
		return ret, nil
	}

	if line[0] == '*' {
		list := bytes.Split(line, []byte("\r\n"))
		nsize, err := strconv.Atoi(string(list[0][1:]))
		if err != nil {
			return make([]byte, 0), err
		}

		rbyte := make([][]byte, nsize)
		var k int = 0
		for i := 1; i < len(list); i++ {
			if bytes.HasPrefix(list[i], []byte("$-1")) {
				k += 1
				continue
			}

			i += 1
			rbyte[k] = unquote(list[i])
			k += 1
		}
		return rbyte, nil
	}

	err = RError("Unkown reply message") // uncatched type
	return make([]byte, 0), err
}

func sendRecv(conn net.Conn, args ...string) (interface{}, error) {
	if conn == nil {
		return nil, RError("Connection is not opened yet!")
	}

	c := strings.Join(args, " ")
	c += "\r\n"
	// fmt.Println(c) // debug

	_, err := conn.Write([]byte(c))
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

func quote(in []byte) []byte {
	out := make([]byte, 0)
	for _, i := range in {
		switch i {
		case '\a':
			out = append(out, `\a`...)
		case '\b':
			out = append(out, `\b`...)
		case '\f':
			out = append(out, `\f`...)
		case '\n':
			out = append(out, `\n`...)
		case '\r':
			out = append(out, `\r`...)
		case '\t':
			out = append(out, `\t`...)
		case '\v':
			out = append(out, `\v`...)
		case 0x0:
			out = append(out, `\0`...)
		case 0x20:
			out = append(out, `\s`...)
		case '\\':
			out = append(out, `\\`...)
		default:
			out = append(out, i)
		}
	}
	return out
}

func unquote(in []byte) []byte {
	out := make([]byte, 0)
	for n := 0; n < len(in); n++ {
		if in[n] == '\\' && n+1 < len(in) {
			switch in[n+1] {
			case 'a':
				out = append(out, '\a')
				n++
			case 'b':
				out = append(out, '\b')
				n++
			case 'f':
				out = append(out, '\f')
				n++
			case 'n':
				out = append(out, '\n')
				n++
			case 'r':
				out = append(out, '\r')
				n++
			case 't':
				out = append(out, '\t')
				n++
			case 'v':
				out = append(out, '\v')
				n++
			case '0':
				out = append(out, 0x0)
				n++
			case 's':
				out = append(out, 0x20)
				n++
			case '\\':
				out = append(out, '\\')
				n++
			default:
				out = append(out, in[n])
			}
		} else {
			out = append(out, in[n])
		}
	}
	return out
	return in
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
