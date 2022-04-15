package generate

import (
	"math/rand"
	"sync"
	"time"
)

const (
	letters              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" //62个 2^6
	letterIdxMaxBits     = 6
	letterIndexMask      = 63 //最多用6为可以表示一个index
	maxRepresentativeNum = 63 / 6
)

var src RandSource = NewRandSource(time.Now().UnixNano())

type RandSource interface {
	Int63() int64
}

type randSource struct {
	mutex  sync.Mutex
	source rand.Source
}

func (rs *randSource) Int63() int64 {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	return rs.source.Int63()
}

func NewRandSource(seed int64) RandSource {
	return &randSource{
		source: rand.NewSource(seed),
		mutex:  sync.Mutex{},
	}
}

func RandStringN(n int) string {

	bytes := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), maxRepresentativeNum; i >= 0; {

		if remain == 0 {
			cache, remain = src.Int63(), maxRepresentativeNum
		}
		if index := int(cache & letterIndexMask); index < len(letters) {
			bytes[i] = letters[index]
			i--
		}
		cache >>= letterIdxMaxBits
		remain--

	}
	return string(bytes)
}
