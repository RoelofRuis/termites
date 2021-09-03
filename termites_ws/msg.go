package termites_ws

import "encoding/json"

const CloseMessageType = "_close"
const ConnectedMessageType = "_connected"
const UpdateMessageType = "update"

func MakeUpdateMessage(data []byte) ([]byte, error) {
	return wrapMessage(UpdateMessageType, data)
}

func MakeCloseMessage() ([]byte, error) {
	return MakeMessage(CloseMessageType, nil)
}

func MakeConnectedMessage(id ClientId) ([]byte, error) {
	data := struct {
		Id ClientId `json:"id"`
	}{
		Id: id,
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
