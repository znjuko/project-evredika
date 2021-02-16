package channels

import (
	"sync"

	"project-evredika/internal/eventer/models"
)

type EventChannel interface {
	Read() chan *models.Event
	Send(task *models.Event)
	Close()
}

type eventCh struct {
	mu     *sync.RWMutex
	closed bool
	c      chan *models.Event
}

func (c *eventCh) Send(task *models.Event) {
	c.mu.RLock()
	defer c.mu.Unlock()
	if !c.closed {
		c.c <- task
	}
}

func (c *eventCh) Read() chan *models.Event {
	return c.c
}

func (c *eventCh) Close() {
	c.mu.RLock()
	defer c.mu.Unlock()
	c.closed = true
	close(c.c)
}

func NewEventChannel(size int) EventChannel {
	return &eventCh{
		c:  make(chan *models.Event, size),
		mu: new(sync.RWMutex),
	}
}
