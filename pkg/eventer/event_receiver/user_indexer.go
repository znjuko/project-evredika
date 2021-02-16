package event_receiver

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"project-evredika/internal/eventer/models"
	v1 "project-evredika/pkg/api/v1"
)

type userPutIndexHandler struct {
	userIndexer userIndexer

	logger *logrus.Logger
}

func (u *userPutIndexHandler) Handle(task *models.Metadata) {
	user := &v1.User{}
	if err := json.Unmarshal(task.Data, user); err != nil {
		u.logger.WithFields(logrus.Fields{"error": err}).Warn("failed to unmarshall user data")
		return
	}
	u.userIndexer.PutUser(user)
}

// NewUserPutIndexHandler ...
func NewUserPutIndexHandler(
	userIndexer userIndexer,

	logger *logrus.Logger,
) Handler {
	return &userPutIndexHandler{
		userIndexer: userIndexer,
		logger:      logger,
	}
}

type userDeleteIndexHandler struct {
	userIndexer userIndexer
}

func (u *userDeleteIndexHandler) Handle(task *models.Metadata) {
	u.userIndexer.DeleteUser(task.Key)
}

// NewUserDeleteIndexHandler ...
func NewUserDeleteIndexHandler(
	userIndexer userIndexer,
) Handler {
	return &userDeleteIndexHandler{
		userIndexer: userIndexer,
	}
}
