package termites

import (
	"io"
)

type GraphOptions func(conf *graphConfig)

func WithConsoleLogger() GraphOptions {
	return func(conf *graphConfig) {
		conf.addConsoleLogger = true
	}
}

func CloseOnShutdown(c io.Closer) GraphOptions {
	return AddEventSubscriber(closeOnTeardown{closer: c})
}

func Named(name string) GraphOptions {
	return func(conf *graphConfig) {
		conf.name = name
	}
}

func WithoutSigtermHandler() GraphOptions {
	return func(conf *graphConfig) {
		conf.withSigtermHandler = false
	}
}

func WithoutRunner() GraphOptions {
	return func(conf *graphConfig) {
		conf.addRunner = false
	}
}

func AddEventSubscriber(sub EventSubscriber) GraphOptions {
	return func(conf *graphConfig) {
		conf.subscribers = append(conf.subscribers, sub)
	}
}
