package termites

import (
	"io"
)

type GraphOptions func(conf *graphConfig)

func WithoutSigtermHandler() GraphOptions {
	return func(conf *graphConfig) {
		conf.withSigtermHandler = false
	}
}

func Named(name string) GraphOptions {
	return func(conf *graphConfig) {
		conf.name = name
	}
}

func WithLogger() GraphOptions {
	return func(conf *graphConfig) {
		conf.addLogger = true
	}
}

func WithoutRunner() GraphOptions {
	return func(conf *graphConfig) {
		conf.addRunner = false
	}
}

func CloseOnShutdown(c io.Closer) GraphOptions {
	return AddEventSubscriber(closeOnShutdown{closer: c})
}

func AddEventSubscriber(sub EventSubscriber) GraphOptions {
	return func(conf *graphConfig) {
		conf.subscribers = append(conf.subscribers, sub)
	}
}
