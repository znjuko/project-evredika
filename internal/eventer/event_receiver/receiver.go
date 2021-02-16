package event_receiver

import (
	"context"
	"runtime"
	"sync"

	"project-evredika/internal/eventer/models"
	handlers "project-evredika/pkg/eventer/event_receiver"
)

type Receiver interface {
	Send(task *models.Metadata)
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type receiver struct {
	in channel

	handlers []handlers.Handler
}

func (s *receiver) Send(task *models.Metadata) {
	s.in.Send(task)
}

func (s *receiver) Start(ctx context.Context, wg *sync.WaitGroup) {
	runtime.Gosched()
	defer wg.Done()

	for {
		select {
		case task := <-s.in.Read():
			{
				for _, h := range s.handlers {
					h.Handle(task)
				}
			}

		case <-ctx.Done():
			s.in.Close()
			return
		}

	}
}

// NewReceiver ...
func NewReceiver(
	in channel,

	handlers []handlers.Handler,
) Receiver {
	return &receiver{
		in:       in,
		handlers: handlers,
	}
}
