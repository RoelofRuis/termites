package termites

import (
	"fmt"
	"log"
)

type logger struct {
	msgChan chan MessageRef
	close   chan bool
}

func newLogger() *logger {
	return &logger{
		msgChan: make(chan MessageRef, 1024),
		close:   make(chan bool),
	}
}

func (l *logger) Setup(registry NodeRegistry) {
	for _, n := range registry.Iterate() {
		n.SetMessageRefChannel(l.msgChan)
	}

	go func() {
		for {
			select {
			case msg := <-l.msgChan:
				if msg.error != nil {
					log.Printf(
						"MESSAGE ERROR: %s\n%s",
						formatRoute(msg),
						msg.error.Error(),
					)
				} else {
					log.Printf(
						"MESSAGE: %s\n",
						formatRoute(msg),
					)
				}
			case <-l.close:
				return
			}
		}
	}()
}

func (l *logger) Teardown() {
	l.close <- true
}

func formatRoute(ref MessageRef) string {
	adapterString := ""
	if ref.adapterName != "" {
		adapterString = fmt.Sprintf("(%s) -> ", ref.adapterName)
	}
	ownerString := ""
	if ref.toName != "" {
		ownerString = fmt.Sprintf("%s:%s", ref.toName, ref.toPortName)
	}
	return fmt.Sprintf("[%s:%s -> %s%s]",
		ref.fromName,
		ref.fromPortName,
		adapterString,
		ownerString,
	)
}
