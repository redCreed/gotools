package queue

import (
	"sync"
	"time"
)

type queue chan Messager

type Memory struct {
	syncMap sync.Map
	wait    sync.WaitGroup
	mutex   sync.RWMutex
	poolNum int
}

func NewMemory(num int) Queue {
	return &Memory{poolNum: num}
}

func (m *Memory) string() string {
	return "memory"
}

func (m *Memory) makeQueue() queue {
	if m.poolNum <= 0 {
		return make(queue)
	}

	return make(queue, m.poolNum)
}

func (m *Memory) Add(messager Messager) error {
	var q queue
	value, ok := m.syncMap.Load(messager.GetKey())
	if !ok {
		q = m.makeQueue()
		m.syncMap.Store(messager.GetKey(), q)
	} else {
		q = value.(queue)
	}

	//推送
	go func(m Messager, q queue) {
		q <- m
	}(messager, q)

	//memoryMessage := new(Message)
	//memoryMessage.SetID(messager.GetID())
	//memoryMessage.SetStream(messager.GetStream())
	//memoryMessage.SetValues(messager.GetValues())
	//v, ok := m.syncMap.Load(messager.GetStream())
	//if !ok {
	//	v = m.makeQueue()
	//	//stream==uuid
	//	m.syncMap.Store(messager.GetStream(), v)
	//}
	//var q queue
	//switch v.(type) {
	//case queue:
	//	q = v.(queue)
	//default:
	//	q = m.makeQueue()
	//	m.syncMap.Store(messager.GetStream(), q)
	//}
	////todo 是否需要协程
	//go func(gm Messager, gq queue) {
	//	gm.SetID(uuid.New().String())
	//	gq <- gm
	//}(memoryMessage, q)
	return nil
}

func (m *Memory) Register(key string, consumerFunc ConsumerFunc) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.syncMap.Load(key)
	if !ok {
		v = m.makeQueue()
		m.syncMap.Store(key, v)
	}
	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.syncMap.Store(key, q)
	}
	go func(out queue, gf ConsumerFunc) {
		var err error
		for message := range out {
			//执行函数
			err = gf(message)
			if err != nil {
				if message.GetErrorCount() < 3 {
					message.SetErrorCount(message.GetErrorCount() + 1)
					// 每次间隔时长放大
					i := time.Second * time.Duration(message.GetErrorCount())
					time.Sleep(i)
					//重新入队
					out <- message
				}
				err = nil
			}
		}
	}(q, consumerFunc)
	return
}

func (m *Memory) Run() {
	m.wait.Add(1)
	m.wait.Wait()
}

func (m *Memory) Close() {
	m.wait.Done()
}
