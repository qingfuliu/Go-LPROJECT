package goRoutinePool

import (
	"fmt"
	"time"
)

type goWorker struct {
	pool        *Pool
	task        chan func()
	recycleTime time.Time
}

func (g *goWorker) run() {
	g.pool.incrRunning()
	go func() {
		defer func() {
			g.pool.decrRunning()
			g.pool.workerCache.Put(g)
			if p := recover(); p != nil {
				fmt.Println(p)
			}
			g.pool.cond.Signal()
		}()
		for f := range g.task {
			if f == nil {
				return
			}
			f()
			if !g.pool.reserveWorkers(g) {
				return
			}
		}
	}()
}
