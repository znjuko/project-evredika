package channels

import (
	"sync"

	"project-evredika/internal/eventer/models"
)

type MetaChannel interface {
	Read() chan *models.Metadata
	Send(task *models.Metadata)
	Close()
}

type metaCh struct {
	mu     *sync.RWMutex
	closed bool
	c      chan *models.Metadata
}

func (c *metaCh) Send(task *models.Metadata) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.closed {
		c.c <- task
	}
}

func (c *metaCh) Read() chan *models.Metadata {
	return c.c
}

func (c *metaCh) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	close(c.c)
}

func NewMetaChannel(size int) MetaChannel {
	return &metaCh{
		c:  make(chan *models.Metadata, size),
		mu: new(sync.RWMutex),
	}
}
