package event_sender

import (
	"context"
	"runtime"
	"sync"

	"project-evredika/internal/eventer/models"
)

type Sender interface {
	SendEvent(events ...*models.Event)
	Start(ctx context.Context, wg *sync.WaitGroup)
	Subscribe(eventType string, subs receiver)
}

type sender struct {
	in channel

	subscribers map[string][]receiver
}

func (s *sender) SendEvent(events ...*models.Event) {
	for _, event := range events {
		s.in.Send(event)
	}
}

func (s *sender) Start(ctx context.Context, wg *sync.WaitGroup) {
	runtime.Gosched()
	defer wg.Done()

	for {
		select {
		case event := <-s.in.Read():
			{
				subs := s.subscribers[event.Type]
				for _, sub := range subs {
					sub.Send(event.Metadata)
				}
			}

		case <-ctx.Done():
			s.in.Close()
			return
		}

	}
}

// Used only on service start
// cannot dynamically add subses (mb i will add it later)
func (s *sender) Subscribe(eventType string, subs receiver) {
	s.subscribers[eventType] = append(s.subscribers[eventType], subs)
}

// NewSender ...
func NewSender(
	in channel,
) Sender {
	return &sender{
		in:          in,
		subscribers: make(map[string][]receiver),
	}
}
