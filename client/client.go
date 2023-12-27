package client

import (
	"fmt"
	"net/http"
)

type Client struct {
	host string
	port string
	user string
	pass string

	url    string
	client *http.Client
}

func New(host, port string) *Client {
	client := &Client{
		host:   host,
		port:   port,
		client: &http.Client{},
	}
	client.updateURL()
	return client
}

func (c *Client) updateURL() {
	var baseAuth string
	if c.user != "" && c.pass != "" {
		baseAuth = fmt.Sprintf("%s:%s@", c.user, c.pass)
	}
	c.url = fmt.Sprintf("http://%s%s:%s/", baseAuth, c.host, c.port)
}
