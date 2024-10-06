package termites_web

import (
	"encoding/json"
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

type WebMessage struct {
	MsgType     string          `json:"msg_type"`
	ContentType string          `json:"content_type"`
	Payload     json.RawMessage `json:"payload"`
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

const (
	MsgClose  = "_close"
	MsgUpdate = "update"
)

func WebClose() ([]byte, error) {
	return marshalWebMessage(WebMessage{MsgType: MsgClose})
}

func WebUpdate(contentType string, data []byte) ([]byte, error) {
	return marshalWebMessage(WebMessage{
		MsgType:     MsgUpdate,
		ContentType: contentType,
		Payload:     data,
	})
}

func marshalWebMessage(msg WebMessage) ([]byte, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgBytes, nil
}
