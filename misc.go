package redis

const (
	poolSize = 10
)

type RedisPool struct {
	Remote string
	Psw    string
	Db     int
	pool   chan *Client
}

func (this *RedisPool) CreatePool() {
	this.pool = make(chan *Client, poolSize)
	for i := 0; i < poolSize; i++ {
		cli := &Client{
			Remote: this.Remote,
			Psw:    this.Psw,
			Db:     this.Db,
		}
		this.pool <- cli
	}
}

func (this *RedisPool) PopClient() (*Client, error) {
	cli := <-this.pool

	if !cli.IsActive() {
		if err := cli.Connect(); err != nil {
			return nil, err
		}
	}
	return cli, nil
}

func (this *RedisPool) PushClient(cli *Client) {
	this.pool <- cli
}

func (this *RedisPool) DestroyPool() {
	this.pool = nil
}
