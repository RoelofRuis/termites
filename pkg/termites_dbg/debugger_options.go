package termites_dbg

import "github.com/gorilla/websocket"

type DebuggerOption func(conf *debuggerConfig)

func OnHttpPort(port int) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.httpPort = port
	}
}

func WithUpgrader(upgrader websocket.Upgrader) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.upgrader = upgrader
	}
}

func WithoutMessageTracking() DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.trackMessages = false
	}
}
