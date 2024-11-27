package historian

import "time"

type Historian interface {
	SaveMessage(msg interface{})
	GetAllMessages() []interface{}
}

type hist struct {
	ttl time.Duration
	ll  *linkedlist
}

func NewHistorian(maxMessages int, ttl time.Duration) Historian {
	return &hist{
		ttl: ttl,
		ll:  newLinkedList(maxMessages, ttl),
	}
}

func (h *hist) SaveMessage(msg interface{}) {
	h.ll.saveValue(msg)
}

func (h *hist) GetAllMessages() []interface{} {
	return h.ll.getAllValues()
}
