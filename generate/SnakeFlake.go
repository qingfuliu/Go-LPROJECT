package generate

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	//机器id
	workbenchIdLength = 4
	workbenchIdShift  = 0
	//机房id
	engineRoomLength = 4
	engineRoomShift  = workbenchIdLength + workbenchIdLength
	//用户id
	idLength = 15
	idShift  = engineRoomLength + engineRoomShift
	maxId    = 1 << idLength
	idMask   = 1<<idLength - 1
	//时间戳
	timeStampLength = 41
	timeStampShift  = idLength + idShift
	timeStampMask   = 1<<timeStampLength - 1
)

var (
	workbenchId int64 = 1
	engineRoom  int64 = 1
	currentId   int64 = 0
	mutex       sync.Mutex
	now         int64
	last        int64 = time.Now().UnixMilli()
)

func GetSnakeId() int64 {
	mutex.Lock()
	defer mutex.Unlock()
	now = time.Now().UnixMilli()
	if now < last {
		zap.L().Fatal("SnakeFlake failed,time is Wrong")
		panic("SnakeFlake failed,time is Wrong")
	}
	if now == last {
		currentId = (currentId + 1) & idMask
		if currentId == 0 {
			nextMilli()
		}
	} else {
		currentId = 0
	}
	last = now
	return workbenchId<<workbenchIdShift |
		engineRoom<<engineRoomShift |
		currentId<<idShift |
		(now&timeStampMask)<<timeStampShift
}

func nextMilli() {
	for now == last {
		now = time.Now().UnixMicro()
	}
	currentId = 0
}
