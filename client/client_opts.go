package client

import "time"

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.client.Timeout = timeout
	return c
}

func (c *Client) WithBaseAuth(user, pass string) *Client {
	c.user = user
	c.pass = pass
	c.updateURL()
	return c
}
