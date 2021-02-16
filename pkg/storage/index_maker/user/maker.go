package user

import "fmt"

type Indexer interface {
	CreateIndex(id string) (index string)
	RevertIndex(index string) (id string)
}

type maker struct {
	format string
	logger logger
}

func (i *maker) RevertIndex(index string) (id string) {
	if _, err := fmt.Sscanf(index, i.format, &id); err != nil {
		i.logger.Warn(
			"failed to revert index", index,
			"format", i.format,
			"error", err,
		)
	}
	return
}

func (i *maker) CreateIndex(id string) (index string) {
	return fmt.Sprintf(i.format, id)
}

// NewIndexer ...
func NewIndexMaker(format string, logger logger) Indexer {
	return &maker{
		format: format,
		logger: logger,
	}
}
