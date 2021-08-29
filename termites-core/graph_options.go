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
	return AddObserver(closeOnShutdown{closer: c})
}

func AddObserver(obs GraphObserver) GraphOptions {
	return func(conf *graphConfig) {
		conf.observers = append(conf.observers, obs)
	}
}
