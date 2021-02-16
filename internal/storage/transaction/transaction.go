package transaction

import (
	"sync"
	"sync/atomic"
)

// Transaction is used to start transaction over storage;
// is used only for one bucket
// it is just a demo version of transactions
type Transaction interface {
	dataSaver
	StartTransaction(key string)
	StopTransaction(key string)
}

type transaction struct {
	dataSaver

	mu *sync.Mutex

	keyTx map[string]*int32

	swapValue    int32
	defaultValue int32
}

func (t *transaction) StartTransaction(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exist := t.keyTx[key]; !exist {
		defVal := t.defaultValue
		t.keyTx[key] = &defVal
	}

	for atomic.SwapInt32(t.keyTx[key], t.swapValue) == t.swapValue {
	}

	return
}

func (t *transaction) StopTransaction(key string) {
	atomic.SwapInt32(t.keyTx[key], t.defaultValue)
	return
}

// NewTransaction ...
func NewTransaction(dataSaver dataSaver, swapValue, defaultValue int32) Transaction {
	return &transaction{
		dataSaver:    dataSaver,
		mu:           &sync.Mutex{},
		keyTx:        make(map[string]*int32),
		swapValue:    swapValue,
		defaultValue: defaultValue,
	}
}
