package other

import "time"

type TimeTrier interface {
	Stop()
	Trier() <-chan time.Time
}

func NewTimeTrier(duration time.Duration) TimeTrier {
	return &timeTrier{
		Ticker: time.NewTicker(duration),
	}
}

type timeTrier struct {
	*time.Ticker
}

func (t *timeTrier) Stop() {
	t.Stop()
}

func (t *timeTrier) Trier() <-chan time.Time {
	return t.C
}
