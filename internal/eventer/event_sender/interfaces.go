package event_sender

import (
	"project-evredika/internal/eventer/models"
)

type channel interface {
	Read() chan *models.Event
	Send(task *models.Event)
	Close()
}

type receiver interface {
	Send(task *models.Metadata)
}
