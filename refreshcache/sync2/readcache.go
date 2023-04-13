package sync2

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type noCopy struct{}

type ReadCache struct {
	noCopy noCopy //nolint:structcheck

	started Fence
	ctx     context.Context

	// read is a func that's called when a new update is needed.
	read func(ctx context.Context) (interface{}, error)
	// refresh defines when the state should be updated.
	refresh time.Duration
	// stale defines when we must wait for the new state.
	stale time.Duration

	// mu protects the internal state of the cache.
	mu sync.Mutex
	// closed is set true when the read cache is shuting down.
	closed bool
	// result contains the last known state and any errors that
	// occurred during refreshing.
	result *readCacheResult
	// pending is a channel for waiting for the current refresh.
	// it is only present, when there is an ongoing refresh.
	pending *readCacheWorker
}

// readCacheResult contains the result of a read and info related to it.
type readCacheResult struct {
	start time.Time
	state interface{}
	err   error
}

// readCacheWorker contains the pending result.
type readCacheWorker struct {
	done   chan struct{}
	result *readCacheResult
}

// NewReadCache returns a new ReadCache.
func NewReadCache(refresh time.Duration, stale time.Duration, read func(ctx context.Context) (interface{}, error)) (*ReadCache, error) {
	cache := &ReadCache{}
	return cache, cache.init(refresh, stale, read)
}

// Init initializes the cache for in-place initialization. This is only needed when NewReadCache
// was not used to initialize it.
func (cache *ReadCache) init(refresh time.Duration, stale time.Duration, read func(ctx context.Context) (interface{}, error)) error {
	if refresh > stale {
		refresh = stale
	}
	if refresh <= 0 || stale <= 0 {
		return errors.Errorf("refresh and stale must be positive. refresh=%v, stale=%v", refresh, stale)
	}
	cache.read = read
	cache.refresh = refresh
	cache.stale = stale
	return nil
}

// Run starts the background process for the cache.
func (cache *ReadCache) Run(ctx context.Context) error {
	// set the root context
	cache.ctx = ctx
	//关闭done
	cache.started.Release()

	// wait for things to start shutting down
	<-ctx.Done()

	// close the workers
	cache.mu.Lock()
	cache.closed = true
	pending := cache.pending
	cache.mu.Unlock()

	// wait for worker to exit
	if pending != nil {
		<-pending.done
	}

	return nil
}

// Get fetches the latest state and refreshes when it's needed.
func (cache *ReadCache) Get(ctx context.Context, now time.Time) (state interface{}, err error) {
	if !cache.started.Wait(ctx) {
		return nil, ctx.Err()
	}

	// check whether we need to start a refresh
	cache.mu.Lock()
	mustWait := false
	//没有缓存数据 有错误  缓存过期
	if cache.result == nil || cache.result.err != nil || now.Sub(cache.result.start) >= cache.refresh {
		// check whether we must wait for the result:
		//   * we don't have anything in cache
		//   * the cache state has errored
		//   * we have reached the staleness deadline
		mustWait = cache.result == nil || cache.result.err != nil || now.Sub(cache.result.start) >= cache.stale
		//刷新缓存
		if err := cache.startRefresh(now); err != nil {
			cache.mu.Unlock()
			return nil, err
		}
	}
	//第一个是缓存中的数据， pending.result则表示正在刷新的数据
	result, pending := cache.result, cache.pending
	//当pending为空，则表示使用的缓存，没有使用刷新数据
	fmt.Printf("refreshing and get pending addr: %p\n", pending)
	cache.mu.Unlock()

	// wait for the new result, when needed
	if mustWait {
		select {
		case <-pending.done:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		result = pending.result
	}

	return result.state, result.err
}

// startRefresh starts a new background refresh, when one isn't running
// already. It will return an error when the cache is shutting down.
//
// Note: this must only be called when `cache.mu` is being held.
func (cache *ReadCache) startRefresh(now time.Time) error {
	if cache.closed {
		return context.Canceled
	}
	if cache.pending != nil {
		return nil
	}

	pending := &readCacheWorker{
		done:   make(chan struct{}),
		result: nil,
	}

	go func() {
		fmt.Println("begin refresh!!! ")
		//外面get函数的时候会阻塞到pending.done中
		defer close(pending.done)

		state, err := cache.read(cache.ctx)
		fmt.Println("use lock begin assignment!!! ")
		cache.mu.Lock()
		result := &readCacheResult{
			start: now,
			state: state,
			err:   err,
		}
		//刷新完直接赋值给result，不需要刷新时可以使用result
		cache.result = result
		//pending已经把该地址传递给调用端，直接把result赋值pending，外部可以使用result
		pending.result = result
		//此次赋值nil表示是刷新缓存结束，但刷新结果已经传递给pending
		cache.pending = nil

		cache.mu.Unlock()
	}()

	fmt.Printf("init pending and addr: %p\n", pending)

	//在外部接收pending地址变量  result, pending := cache.result, cache.pending
	cache.pending = pending
	return nil
}
