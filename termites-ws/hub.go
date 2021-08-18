package cliserv

import (
	"sync"

	"github.com/google/uuid"
)

type Hub interface {
	Broadcast(data []byte)
	Send(data []byte, receivers ...ClientId)
	ReadReceive() chan ClientMessage
	ReadConnect() chan ClientConnection
}

type ClientMessage struct {
	Sender ClientId
	Data   []byte
}

type ClientConnection struct {
	ConnType ConnectionType
	Id       ClientId
}

type ConnectionType uint8

const (
	ClientConnect    ConnectionType = 0
	ClientDisconnect ConnectionType = 1
	ClientReconnect  ConnectionType = 2
)

type MessageHub struct {
	clientsLock   *sync.RWMutex
	clients       map[ClientId]*Client
	receiveIn     chan ClientMessage
	receiveOut    chan ClientMessage
	connectionOut chan ClientConnection
	isReceiving   bool
	register      chan *Client
	unregister    chan *Client
}

func NewHub() *MessageHub {
	hub := &MessageHub{
		clientsLock:   &sync.RWMutex{},
		clients:       make(map[ClientId]*Client),
		receiveIn:     make(chan ClientMessage),
		receiveOut:    make(chan ClientMessage),
		connectionOut: make(chan ClientConnection),
		isReceiving:   false,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}

	go func() {
		for {
			select {
			case client := <-hub.register:
				hub.clientsLock.Lock()
				_, has := hub.clients[client.id]
				hub.clients[client.id] = client
				hub.clientsLock.Unlock()
				connType := ClientConnect
				if has {
					connType = ClientReconnect
				}
				hub.connectionOut <- ClientConnection{ConnType: connType, Id: client.id}

				if connMsg, err := MakeConnectedMessage(client.id); err == nil {
					// send if no error
					client.Received <- connMsg
				}

			case client := <-hub.unregister:
				hub.clientsLock.Lock()
				c, has := hub.clients[client.id]
				if has && c != nil {
					hub.clients[client.id] = nil
					close(c.Received)
				}
				hub.clientsLock.Unlock()
				hub.connectionOut <- ClientConnection{ConnType: ClientDisconnect, Id: client.id}

			case msg := <-hub.receiveIn:
				if !hub.isReceiving {
					continue
				}

				hub.receiveOut <- msg
			}
		}
	}()

	return hub
}

func (e *MessageHub) RegisterClient(id string) *Client {
	var clientId = ClientId(uuid.NewString())

	e.clientsLock.RLock()
	if c, has := e.clients[ClientId(id)]; has && c == nil {
		clientId = ClientId(id)
	}
	e.clientsLock.RUnlock()

	client := &Client{
		id:         clientId,
		receive:    e.receiveIn,
		unregister: e.unregister,
		Received:   make(chan []byte, 256),
	}

	e.register <- client

	return client
}

func (e *MessageHub) Broadcast(data []byte) {
	e.clientsLock.RLock()
	defer e.clientsLock.RUnlock()

	for _, c := range e.clients {
		if c == nil {
			continue
		}
		c.Received <- data
	}
}

func (e *MessageHub) Send(data []byte, receivers ...ClientId) {
	e.clientsLock.RLock()
	defer e.clientsLock.RUnlock()

	for _, id := range receivers {
		c, has := e.clients[id]
		if !has || c == nil {
			continue
		}
		c.Received <- data
	}
}

func (e *MessageHub) ReadReceive() chan ClientMessage {
	e.isReceiving = true
	return e.receiveOut
}

func (e *MessageHub) ReadConnect() chan ClientConnection {
	return e.connectionOut
}
