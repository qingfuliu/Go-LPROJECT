package goRoutinePool

import (
	"math"
	"time"
)

const (
	poolStateRunning = iota
	poolStateClosed
	DefaultCleanIntervalTime = time.Second
	DefaultAntsPoolSize      = math.MaxInt32
)

type Option func(options *Options)

func loadOption(options ...Option) *Options {
	opts := &Options{}
	for _, val := range options {
		val(opts)
	}
	return opts
}

type Options struct {
	NoBlocking bool

	MaxBlockNums int32

	ExpireDuration time.Duration

	PreAllocate bool
}

func WithNoBlocking(noBlocking bool) Option {
	return func(options *Options) {
		options.NoBlocking = noBlocking
	}
}

func WithMaxBlockNums(maxBlockNums int32) Option {
	return func(options *Options) {
		options.MaxBlockNums = maxBlockNums
	}
}

func WithExpireDuration(expireDuration time.Duration) Option {
	return func(options *Options) {
		options.ExpireDuration = expireDuration
	}
}

func WithPreAllocate(preAllocate bool) Option {
	return func(options *Options) {
		options.PreAllocate = preAllocate
	}
}
