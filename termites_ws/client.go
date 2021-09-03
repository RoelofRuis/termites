package termites_ws

type ClientRegistry interface {
	RegisterClient(id string) *Client
}

type ClientId string

type Client struct {
	id         ClientId
	receive    chan ClientMessage
	unregister chan *Client
	Received   chan []byte
}

func (c *Client) Send(data []byte) {
	c.receive <- ClientMessage{Sender: c.id, Data: data}
}

func (c *Client) Unregister() {
	c.unregister <- c
}
