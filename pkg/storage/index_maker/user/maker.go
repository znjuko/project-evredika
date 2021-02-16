package user

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Indexer interface {
	CreateIndex(id string) (index string)
	RevertIndex(index string) (id string)
}

type maker struct {
	suffix string
	logger *logrus.Logger
}

func (i *maker) RevertIndex(index string) (id string) {
	id = strings.TrimSuffix(index, i.suffix)
	i.logger.WithFields(
		logrus.Fields{
			"index":  index,
			"suffix": i.suffix,
			"id":     id,
		}).Debug("RevertIndex")
	return
}

func (i *maker) CreateIndex(id string) (index string) {
	index = id + i.suffix

	i.logger.WithFields(
		logrus.Fields{
			"index":  index,
			"suffix": i.suffix,
			"id":     id,
		}).Debug("CreateIndex")

	return
}

// NewIndexer ...
func NewIndexMaker(suffix string, logger *logrus.Logger) Indexer {
	return &maker{
		suffix: suffix,
		logger: logger,
	}
}
