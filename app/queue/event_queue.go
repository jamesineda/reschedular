package queue

import (
	"github.com/jamesineda/reschedular/app/event"
	"sync"
)

type Events struct {
	sync.Mutex
	Queue event.IncomingEvents
}

// NewEventsQueue setting a capacity just alleviates some of the re-sizing capacity
func NewEventsQueue(capacity int) *Events {
	return &Events{
		Queue: make(event.IncomingEvents, 0, capacity),
	}
}

func (q *Events) Pop() event.IncomingEvent {
	q.Lock()
	defer q.Unlock()
	if len(q.Queue) == 0 {
		return nil
	}

	e := q.Queue[0]
	q.Queue[0] = nil
	q.Queue = q.Queue[1:]
	return e
}

func (q *Events) Push(instruction event.IncomingEvent) {
	q.Lock()
	defer q.Unlock()
	q.Queue = append(q.Queue, instruction)
}
