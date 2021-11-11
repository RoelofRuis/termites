package termites_dbg

type DebuggerOption func(conf *debuggerConfig)

func OnHttpPort(port int) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.httpPort = port
	}
}

func OpenIn(editor CodeEditor) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.editor = editor
	}
}
