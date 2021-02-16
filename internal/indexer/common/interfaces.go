package common

type lister interface {
	Add(key string)
	Remove(key string)
	List(skip, limit int) (keys []string)
}
