package termites

import (
	"io"
)

type GraphOptions func(conf *graphConfig)

func WithoutSigtermHandler() GraphOptions {
	return func(conf *graphConfig) {
		conf.withSigtermHandler = true
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

func NonblockingRun() GraphOptions {
	return func(conf *graphConfig) {
		conf.blockingRun = false
	}
}

func CloseOnShutdown(c io.Closer) GraphOptions {
	return AddHook(closeOnShutdown{closer: c})
}

func AddHook(hook GraphHook) GraphOptions {
	return func(conf *graphConfig) {
		conf.hooks = append(conf.hooks, hook)
	}
}
