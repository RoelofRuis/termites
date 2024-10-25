package termites

import (
	"io"
)

// graphConfig holds configuration data to be set via various GraphOption implementations.
type graphConfig struct {
	name               string
	subscribers        []EventSubscriber
	withSigtermHandler bool
	printLogs          bool
	printMessages      bool
	deferredStart      bool
}

// GraphOption allows modifications to the graphConfig
type GraphOption func(conf *graphConfig)

func PrintLogsToConsole() GraphOption {
	return func(conf *graphConfig) {
		conf.printLogs = true
	}
}

func PrintMessagesToConsole() GraphOption {
	return func(conf *graphConfig) {
		conf.printMessages = true
	}
}

func CloseOnTeardown(resourceName string, c io.Closer) GraphOption {
	return WithEventSubscriber(closeOnTeardown{name: resourceName, closer: c})
}

func Named(name string) GraphOption {
	return func(conf *graphConfig) {
		conf.name = name
	}
}

func WithoutSigtermHandler() GraphOption {
	return func(conf *graphConfig) {
		conf.withSigtermHandler = false
	}
}

func WithEventSubscriber(sub EventSubscriber) GraphOption {
	return func(conf *graphConfig) {
		conf.subscribers = append(conf.subscribers, sub)
	}
}

// DeferredStart defers running the nodes until Graph.Wait is called.
func DeferredStart() GraphOption {
	return func(conf *graphConfig) {
		conf.deferredStart = true
	}
}
