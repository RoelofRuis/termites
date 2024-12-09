package termites_dbg

import (
	"github.com/RoelofRuis/termites/pkg/termites"
	"github.com/RoelofRuis/termites/pkg/termites_web"
)

var MessageSentAdapter = func(event termites.MessageSentEvent) (termites_web.ClientMessage, error) {
	return termites_web.NewClientMessage("message", event)
}

var LogsAdapter = func(event logItem) (termites_web.ClientMessage, error) {
	return termites_web.NewClientMessage("log", event)
}
