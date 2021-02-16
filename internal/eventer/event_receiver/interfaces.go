package event_receiver

import "project-evredika/internal/eventer/models"

type channel interface {
	Read() chan *models.Metadata
	Send(task *models.Metadata)
	Close()
}
