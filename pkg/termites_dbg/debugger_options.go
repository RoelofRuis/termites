package termites_dbg

import "github.com/gorilla/websocket"

type debuggerConfig struct {
	httpPort int
	upgrader websocket.Upgrader

	trackRefChanges bool
	trackMessages   bool
	trackLogs       bool
}

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

func WithoutRefTracking() DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.trackRefChanges = true
	}
}

func WithoutLogTracking() DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.trackLogs = false
	}
}
