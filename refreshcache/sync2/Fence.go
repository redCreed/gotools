package sync2

import (
	"context"
	"sync"
)

// Fence allows to wait for something to happen.
type Fence struct {
	noCopy  noCopy
	setup   sync.Once
	release sync.Once
	done    chan struct{}
}

// init sets up the initial lock into wait.
func (f *Fence) init() {
	f.setup.Do(func() {
		f.done = make(chan struct{})
	})
}

// Release releases everyone from Wait.
func (f *Fence) Release() {
	f.release.Do(func() {
		close(f.done)
	})
}

// Wait waits for wait to be unlocked.
// Returns true when it was successfully released.
func (f *Fence) Wait(ctx context.Context) bool {
	f.init()
	select {
	case <-f.done:
		return true
	default:
		select {
		case <-ctx.Done():
			return false
		case <-f.done:
			return true
		}
	}
}

// Released returns whether the fence has been released.
func (f *Fence) Released() bool {
	f.init()

	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// Done returns channel that will be closed when the fence is released.
func (f *Fence) Done() chan struct{} {
	f.init()
	return f.done
}
