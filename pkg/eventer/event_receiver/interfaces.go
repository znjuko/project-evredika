package event_receiver

import v1 "project-evredika/pkg/api/v1"

type userIndexer interface {
	PutUser(user *v1.User)
	DeleteUser(ID string)
}

type logger interface {
	Warn(args ...interface{})
}
