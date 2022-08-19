package paperId

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1629167031000
)

type SnowWorker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewSnowWorker(workerId int64) (*SnowWorker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	return &SnowWorker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *SnowWorker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	return (now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number)
}
