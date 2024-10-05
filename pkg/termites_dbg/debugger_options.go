package termites_dbg

type DebuggerOption func(conf *debuggerConfig)

func OnHttpPort(port int) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.httpPort = port
	}
}

// Deprecated
// Determining function file and line was broken in go 1.18
func OpenIn(editor CodeEditor) DebuggerOption {
	return func(conf *debuggerConfig) {
		conf.editor = editor
	}
}
