package termites_dbg

import (
	"github.com/RoelofRuis/termites/pkg/termites"
)

type logReceiver struct {
	LogsOut *termites.OutPort
}

func newLogReceiver() *logReceiver {
	builder := termites.NewBuilder("Log Receiver")

	n := &logReceiver{
		LogsOut: termites.NewOutPortNamed[logItem](builder, "Logs"),
	}

	return n
}

func (l *logReceiver) onLog(logEvent termites.Event) error {
	switch evt := logEvent.Data.(type) {
	case termites.InfoLogEvent:
		l.LogsOut.Send(logItem{LogLevel: "info", Message: evt.Message})
		break

	case termites.ErrorLogEvent:
		l.LogsOut.Send(logItem{LogLevel: "error", Message: evt.Message, Error: evt.Error.Error()})
		break

	case termites.PanicLogEvent:
		l.LogsOut.Send(logItem{LogLevel: "panic", Message: evt.Message, Stack: evt.Stack})

	default:
		return termites.InvalidEventError
	}

	return nil
}

type logItem struct {
	LogLevel string `json:"log_level"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Stack    string `json:"stack_trace"`
}
