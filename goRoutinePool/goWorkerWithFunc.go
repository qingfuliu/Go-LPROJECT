package goRoutinePool

import (
	"fmt"
	"time"
)

type goWorkerWithFunc struct {
	pool        *PoolFunc
	recycleTime time.Time
	arg         chan interface{}
}

func (g *goWorkerWithFunc) run() {
	g.pool.incRunning()
	defer func() {
		g.pool.decrRunning()
		g.pool.reserveWorkers(g)
		if p := recover(); p != nil {
			//logg
			fmt.Println(p)
		}
		g.pool.cond.Signal()
	}()
	for arg := range g.arg {
		if arg == nil {
			return
		}
		g.pool.goFunc(arg)
		if !g.pool.reserveWorkers(g) {
			return
		}
	}
}
