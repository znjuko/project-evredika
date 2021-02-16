package common

import (
	"sync"

	v1 "project-evredika/pkg/api/v1"
)

type UserIndexer interface {
	PutUser(user *v1.User)
	DeleteUser(ID string)
	GetUser(ID string) (user *v1.User, exist bool)
	ListUsers(skip, limit int) (users []*v1.User)
}

type userIndexer struct {
	mu *sync.RWMutex

	users map[string]*v1.User

	lister lister
}

func (i *userIndexer) PutUser(user *v1.User) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.lister.Add(user.ID)
	i.users[user.ID] = user
	return
}

func (i *userIndexer) DeleteUser(ID string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.lister.Remove(ID)
	if _, exist := i.users[ID]; exist {
		delete(i.users, ID)
	}
	return
}

func (i *userIndexer) GetUser(ID string) (user *v1.User, exist bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	user, exist = i.users[ID]
	return
}

func (i *userIndexer) ListUsers(skip, limit int) (users []*v1.User) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	idxs := i.lister.List(skip, limit)
	for _, idx := range idxs {
		if user, exist := i.users[idx]; exist {
			users = append(users, user)
		}
	}

	return
}

// NewUserIndexer ...
func NewUserIndexer(lister lister) UserIndexer {
	return &userIndexer{
		mu:     &sync.RWMutex{},
		users:  make(map[string]*v1.User),
		lister: lister,
	}
}
