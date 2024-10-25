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

	// Data any data to be sent to the client.
	Data []byte
}

// NewClientMessage creates a new client message with the given topic and any data.
// If the data is already of type []byte, it will be sent as is, otherwise, json encoding will be performed first.
func NewClientMessage(topic string, data any) (ClientMessage, error) {
	return NewClientMessageFor(topic, "", data)
}

func NewClientMessageFor(topic string, clientId string, data any) (message ClientMessage, err error) {
	payload, is := data.([]byte)
	if !is {
		payload, err = json.Marshal(data)
		if err != nil {
			return message, err
		}
	}

	clientData, err := json.Marshal(WebMessage{
		Topic:   topic,
		Payload: payload,
	})
	if err != nil {
		return message, err
	}

	message.ClientId = clientId
	message.Data = clientData
	return message, nil
}

const SystemCloseTopic = "system/close"

type WebMessage struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// WebMsgAdapter unpacks the ClientMessage data into a WebMessage.
func WebMsgAdapter(c ClientMessage) (WebMessage, error) {
	msg := WebMessage{}
	if err := json.Unmarshal(c.Data, &msg); err != nil {
		return msg, err
	}
	return msg, nil
}
