package event_receiver

import "project-evredika/internal/eventer/models"

type Handler interface {
	Handle(task *models.Metadata)
}
