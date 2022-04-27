package goRoutinePool

import "errors"

var (
	errorQueueIsFull     = errors.New(" Queue Is Full")
	errorInvalidExpire   = errors.New("expire is invalid")
	errorInvalidCapacity = errors.New("capacity is invalid")
	errorPoolOverLoad    = errors.New(" Pool Is overload")
	errorPoolClosed      = errors.New("pool is closed")
)
