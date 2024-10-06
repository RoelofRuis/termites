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

// Deprecated
// Determining function file and line was broken in go 1.18
func OpenIn(editor CodeEditor) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.editor = editor
	}
}
