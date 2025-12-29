package main

type MatLock struct {
	locked bool
	wait chan struct{}
}

func NewQutex() *MatLock {
	return &MatLock{
        locked: false,
        wait:   make(chan struct{}, 1),
	}
}

func (q *MatLock) Lock() {
	for {
     	if !q.locked {
          	q.locked = true
          	return
     	}
		<-q.wait
	}
}

func (q *MatLock) Unlock() {
	if !q.locked {
		panic("unlock of unlocked qutex")
	}
	q.locked = false
    
	select {
	case q.wait <- struct{}{}:
	default:
	}
}
