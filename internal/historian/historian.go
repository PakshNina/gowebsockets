package historian

import (
	"sync"
	"time"
)

type Historian interface {
	AddMessage(msg interface{})
	GetAllMessages() []interface{}
	Start()
}

type history struct {
	ll *linkedlist
}

func NewHistorian(maxMessageStored int, ttl time.Duration) Historian {
	return &history{
		ll: &linkedlist{
			currentLength: 0,
			maxLength:     maxMessageStored,
			ttl:           ttl,
			m:             &sync.Mutex{},
			head:          nil,
			tail:          nil,
		},
	}
}

func (h *history) Start() {
	for {
		select {
		case <-time.After(h.ll.ttl / 2):
			h.ll.checkForTTL()
		}
	}
}

func (h *history) AddMessage(msg interface{}) {
	h.ll.saveValue(msg)
}

func (h *history) GetAllMessages() []interface{} {
	return h.ll.getAllMsg()
}
