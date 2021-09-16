package termites

import (
	"io"
)

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
