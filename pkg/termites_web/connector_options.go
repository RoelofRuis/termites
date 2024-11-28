package termites_web

import "github.com/gorilla/websocket"

type connectorConfig struct {
	upgrader websocket.Upgrader
	// Websocket read limit in bytes
	readLimit int64
}

type ConnectorOption func(conf *connectorConfig)

func WithUpgrader(upgrader websocket.Upgrader) ConnectorOption {
	return func(conf *connectorConfig) {
		conf.upgrader = upgrader
	}
}

func WithReadLimit(limitBytes int64) ConnectorOption {
	return func(conf *connectorConfig) {
		conf.readLimit = limitBytes
	}
}
