package termites_web

import (
	"encoding/json"
	"github.com/RoelofRuis/termites/pkg/termites"
)

// ClientMessage is a message used to interact with a web client.
type ClientMessage struct {
	// ClientId is a unique identifier for an attached client.
	// If this field is set on an incoming message, it means the message was sent over websocket by the given client.
	// If this field is set on an outgoing message, it will only be sent to the client with that ID.
	ClientId string

	// Data any data bytes to be sent to the client.
	Data []byte
}

type ClientConnection struct {
	ConnType ConnectionType
	Id       string
}

type ConnectionType uint8

const (
	ClientConnect    ConnectionType = 0
	ClientDisconnect ConnectionType = 1
	ClientReconnect  ConnectionType = 2
)

const CloseMessageType = "_close"
const ConnectedMessageType = "_connected"
const UpdateMessageType = "update"

func MakeUpdateMessage(data []byte) ([]byte, error) {
	return wrapMessage(UpdateMessageType, data)
}

func MakeCloseMessage() ([]byte, error) {
	return MakeMessage(CloseMessageType, nil)
}

func MakeConnectedMessage(id termites.Identifier) ([]byte, error) {
	data := struct {
		Id string `json:"id"`
	}{
		Id: id.Id,
	}
	return MakeMessage(ConnectedMessageType, data)
}

func MakeMessage(tpe string, data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return wrapMessage(tpe, dataBytes)
}

func wrapMessage(tpe string, data []byte) ([]byte, error) {
	msg := message{
		Type: tpe,
		Data: data,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}

type message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
