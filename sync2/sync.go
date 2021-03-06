package sync2

import (
	"sync"
	"sync/atomic"
)

// Once is an object that will perform exactly one action unless call Undo.
type Once struct {
	m    sync.Mutex
	done uint32
}

// Do will do f only once no matter it's successful or not
// if f is blocked, Do will also be
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	if o.done == 0 {
		f()
		atomic.StoreUint32(&o.done, 1)
	}
	o.m.Unlock()
}

// ErrorDo will do f only once when function call is successful,
// if function return error, Once will be failed
func (o *Once) ErrorDo(f func() error) (err error) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	if o.done == 0 {
		if err = f(); err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	o.m.Unlock()
	return
}

// Undo restore Once's state to initial,
// then you can call Do or ErrorDo again
func (o *Once) Undo() {
	o.m.Lock()
	atomic.StoreUint32(&o.done, 0)
	o.m.Unlock()
}

// LockDo do function in lock
func LockDo(lock sync.Locker, fn func()) {
	lock.Lock()
	fn()
	lock.Unlock()
}
