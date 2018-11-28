package workerman_client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/puper/ppgo/helpers"
)

type Client struct {
	Addrs  []string
	config *Config
}

func New(config *Config) (*Client, error) {
	if config.Addr == "" {
		return nil, fmt.Errorf("rpc addr can not be empty")
	}
	return &Client{
		config: config,
		Addrs:  strings.Split(config.Addr, ";"),
	}, nil
}

func (this *Client) getAddr() string {
	return this.Addrs[helpers.GlobalRand().Intn(len(this.Addrs))]
}

func (this *Client) Call(class string, method string, params []interface{}, reply interface{}) error {
	conn, err := net.Dial("tcp", this.getAddr())
	if err != nil {
		return err
	}
	defer conn.Close()
	if this.config.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.config.Timeout)))
	}
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	decoder.UseNumber()
	req := map[string]interface{}{
		"class":       class,
		"method":      method,
		"param_array": params,
	}
	err = encoder.Encode(req)
	if err != nil {
		return err
	}
	resp := &Response{}
	resp.ResponseError = &ResponseError{}
	resp.Data = reply
	err = decoder.Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Code != this.config.SuccessCode {
		return resp.ResponseError
	}
	return nil
}
