package lister

type Lister interface {
	Add(key string)
	Remove(key string)
	List(skip, limit int) (keys []string)
}

type lister struct {
	keys  []string
	exist map[string]struct{}
}

func (l *lister) Add(key string) {
	if _, exist := l.exist[key]; exist {
		return
	}
	l.exist[key] = struct{}{}

	l.keys = append(l.keys, key)
	return
}

func (l *lister) Remove(key string) {
	if _, exist := l.exist[key]; !exist {
		return
	}
	delete(l.exist, key)

	for i, k := range l.keys {
		if k == key {
			l.keys = append(l.keys[:i], l.keys[i+1:]...)
			break
		}
	}
	return
}

func (l *lister) List(skip, limit int) (keys []string) {
	keysLength := len(l.keys)
	if skip >= keysLength {
		return
	}

	limit += skip
	if skip+limit > keysLength {
		limit = keysLength
	}

	return l.keys[skip:limit]
}

// NewLister ...
func NewLister() Lister {
	return &lister{
		keys:  make([]string, 20),
		exist: make(map[string]struct{}, 20),
	}
}
