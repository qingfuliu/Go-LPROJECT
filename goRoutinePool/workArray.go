package goRoutinePool

import "time"

type arrayType int

const (
	stackWorkArray arrayType = 1 << iota
	queueWorkArray
)

type workArray interface {
	insert(g *goWorker) error
	detach() *goWorker
	len() int
	isEmpty() bool
	retrieveExpire(duration time.Duration) []*goWorker
	reset()
}

func newWorkArrayWithType(workArrayType arrayType, size int) workArray {
	switch workArrayType {
	case stackWorkArray:
		return newWorkerStack(size)
	case queueWorkArray:
		return newWorkerQueue(size)
	default:
		return newWorkerStack(size)
	}
}
