package other

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrorShutdown   = errors.New("task does not existed")
	ErrorExisted    = errors.New("task has been existed")
	ErrorIsNil      = errors.New("Unexpected nil value")
	ErrorNotExisted = errors.New("task does not existed")
)

type TimeWheel interface {
	RemoveTask(key interface{}) error
	AddTask(key, value interface{}, delay time.Duration, repeat bool) error
	ModifyTask(key interface{}, delay time.Duration, repeat bool) error
	DrainTask(task func(key, value interface{})) error
	Stop()
	Run()
}

func NewTimeWheel(interval time.Duration, execute Execute, numSlots int, context context.Context) TimeWheel {
	return &timeWheel{
		interval:       interval,
		trier:          NewTimeTrier(interval),
		execute:        execute,
		numSlots:       numSlots,
		context:        context,
		setTaskChan:    make(chan baseEntry),
		removeTaskChan: make(chan interface{}),
		stopChan:       make(chan struct{}),
		modifyTaskChan: make(chan baseEntry),
		drainChan:      make(chan func(key, value interface{})),
		pos:            0,
		tasksMap:       sync.Map{},
		tasksList:      make([]*list.List, numSlots),
	}
}

type Execute func(key, value interface{})
type timeWheel struct {
	tasksMap       sync.Map
	tasksList      []*list.List
	execute        Execute
	trier          TimeTrier
	setTaskChan    chan baseEntry
	removeTaskChan chan interface{}
	stopChan       chan struct{}
	modifyTaskChan chan baseEntry
	trierChan      <-chan time.Time
	drainChan      chan func(key, value interface{})
	pos            int
	numSlots       int
	interval       time.Duration
	context        context.Context
}

type timeingEntry struct {
	*baseEntry
	diff    int
	circle  int
	removed bool
}

type baseEntry struct {
	delay  time.Duration
	key    interface{}
	value  interface{}
	repeat bool
}

type positionEntry struct {
	pos int
	*timeingEntry
}

type TimeTask struct {
	key   interface{}
	value interface{}
}

func (t *timeWheel) RemoveTask(key interface{}) error {
	if key == nil {
		return ErrorIsNil
	}
	select {
	case t.removeTaskChan <- key:
		return nil
	case <-t.stopChan:
		return ErrorShutdown
	}
}
func (t *timeWheel) AddTask(key, value interface{}, delay time.Duration, repeat bool) error {
	select {
	case t.setTaskChan <- baseEntry{
		delay:  delay,
		repeat: repeat,
		key:    key,
		value:  value,
	}:
		return nil
	case <-t.stopChan:
		return ErrorShutdown
	}
}
func (t *timeWheel) ModifyTask(key interface{}, delay time.Duration, repeat bool) error {
	select {
	case t.modifyTaskChan <- baseEntry{
		key:    key,
		repeat: repeat,
		delay:  delay,
	}:
		return nil
	case <-t.stopChan:
		return ErrorShutdown
	}
}
func (t *timeWheel) DrainTask(task func(key, value interface{})) error {
	select {
	case t.drainChan <- task:
		return nil
	case <-t.stopChan:
		return ErrorShutdown
	}
}
func (t *timeWheel) Stop() {
	close(t.stopChan)
}
func (t *timeWheel) Run() {
	t.trierChan = t.trier.Trier()
	for index := 0; index < t.numSlots; index++ {
		t.tasksList[index] = list.New()
	}
	go func() {
		for {
			select {
			case <-t.context.Done():
				t.trier.Stop()
				return
			case <-t.stopChan:
				t.trier.Stop()
				return
			case key := <-t.removeTaskChan:
				t.removeTask(key)
			case base := <-t.setTaskChan:
				t.addTask(&base)
			case base := <-t.modifyTaskChan:
				t.modifyTask(&base)
			case fun := <-t.drainChan:
				t.drainTask(fun)
			case <-t.trierChan:
				t.scanTaskAndRun()
			}
		}
	}()
}

func (t *timeWheel) cacluatePosition(duration time.Duration) (position int, circle int) {
	if duration < t.interval {
		return t.pos + 1, 0
	}
	temp := int(duration / t.interval)
	position = (temp + t.pos) % t.numSlots
	circle = (temp - 1) / t.numSlots
	return
}

func (t *timeWheel) removeTask(key interface{}) {
	value, ok := t.tasksMap.Load(key)
	if !ok {
		return
	}
	task := value.(*positionEntry)
	task.removed = true
	t.tasksMap.Delete(key)
}
func (t *timeWheel) addTask(base *baseEntry) {
	if _, ok := t.tasksMap.Load(base.key); ok {
		return
	}
	if base.delay < t.interval {
		t.execute(base.key, base.value)
		if !base.repeat {
			return
		}
	}
	position, circle := t.cacluatePosition(base.delay)
	timeEntry := &timeingEntry{
		baseEntry: base,
		circle:    circle,
	}
	t.tasksList[position].PushBack(timeEntry)
	t.setOrLoadMap(position, timeEntry)
}
func (t *timeWheel) modifyTask(base *baseEntry) {
	val, ok := t.tasksMap.Load(base.key)
	if !ok {
		return
	}

	posEntry := val.(*positionEntry)
	pos := posEntry.pos

	if base.delay < t.interval {
		posEntry.baseEntry.delay = base.delay
		posEntry.baseEntry.repeat = base.repeat
		return
	}

	position, circle := t.cacluatePosition(base.delay)
	if pos <= position {
		posEntry.pos, posEntry.circle = position, circle
		posEntry.diff = position - pos
	} else if circle > 1 {
		posEntry.pos, posEntry.circle = position, circle-1
		posEntry.diff = t.numSlots + position - pos
	} else {
		posEntry.removed = true
		base.value = posEntry.value
		newTimingEntry := &timeingEntry{
			baseEntry: base,
			circle:    circle,
		}
		t.tasksList[position].PushBack(newTimingEntry)
		t.tasksMap.Store(position, newTimingEntry)
	}
	return
}
func (t *timeWheel) drainTask(task func(key, value interface{})) {
	t.tasksMap.Range(func(key, value interface{}) bool {
		task(key, value)
		val := value.(*positionEntry)
		val.removed = true
		if !val.repeat {
			t.tasksMap.Delete(key)
		} else {
			pos, circle := t.cacluatePosition(val.delay)
			newtimeEntry := &timeingEntry{
				baseEntry: val.baseEntry,
				circle:    circle,
			}
			val.timeingEntry = newtimeEntry
			t.setOrLoadMap(pos, newtimeEntry)
		}
		return true
	})
}

func (t *timeWheel) scanTaskAndRun() {
	t.pos = (t.pos + 1) % t.numSlots
	timeTasks := make([]TimeTask, 0)
	taskList := t.tasksList[t.pos]
	for task := taskList.Front(); task != nil; {
		next := task.Next()
		val := task.Value.(*timeingEntry)
		if val.removed {
			taskList.Remove(task)
			task = next
			continue
		} else if val.diff != 0 {
			taskList.Remove(task)
			newPos:=(t.pos+val.diff)%t.numSlots
			val.diff=0
			t.tasksList[newPos].PushBack(newPos)
			t.setOrLoadMap(newPos,val)
			task = next
			continue
		} else if val.circle != 0 {
			val.circle--
			task = next
			continue
		}

		timeTasks = append(timeTasks, TimeTask{
			value: val.value,
			key:   val.key,
		})
		if val.repeat {
			val.removed = true
			pos, circle := t.cacluatePosition(val.delay)
			newTimeEntry := &timeingEntry{
				baseEntry: val.baseEntry,
				circle:    circle,
			}
			t.tasksList[pos].PushBack(newTimeEntry)
			t.setOrLoadMap(pos, val)
		} else {
			taskList.Remove(task)
		}
		task = next
	}
	for index, _ := range timeTasks {
		t.execute(timeTasks[index].key, timeTasks[index].value)
	}
}

func (t *timeWheel) setOrLoadMap(pos int, time *timeingEntry) {
	if val, ok := t.tasksMap.Load(time.key); ok {
		posEntry := val.(*positionEntry)
		posEntry.pos = pos
		posEntry.timeingEntry = time
	} else {
		t.tasksMap.Store(time.key, &positionEntry{
			pos:          pos,
			timeingEntry: time,
		})
	}
}
