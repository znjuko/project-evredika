package models

const (
	EventDelete = "Delete"
	EventPut = "Put"
)

type Event struct {
	Type string
	*Metadata
}

type Metadata struct {
	Data []byte
	Key  string
}
