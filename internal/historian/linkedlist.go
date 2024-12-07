package historian

import (
	"sync"
	"time"
)

type linkedlist struct {
	currentLength int
	maxLength     int
	ttl           time.Duration
	m             *sync.Mutex
	head          *node
	tail          *node
}

type node struct {
	value     interface{}
	createdAt time.Time
	next      *node
}

func (l *linkedlist) saveValue(v interface{}) {
	l.checkForTTL()
	newNode := &node{
		value:     v,
		createdAt: time.Now(),
		next:      nil,
	}
	l.m.Lock()
	defer l.m.Unlock()
	if l.head == nil {
		l.head = newNode
		l.tail = l.head
		l.currentLength++
		return
	}
	if l.currentLength+1 > l.maxLength {
		l.head = l.head.next
	} else {
		l.currentLength++
	}
	l.tail.next = newNode
	l.tail = l.tail.next
}

func (l *linkedlist) getAllMsg() []interface{} {
	var msgs []interface{}
	l.m.Lock()
	defer l.m.Unlock()
	current := l.head
	for current != nil {
		msgs = append(msgs, current.value)
		current = current.next
	}
	return msgs
}

func (l *linkedlist) checkForTTL() {
	l.m.Lock()
	defer l.m.Unlock()
	current := l.head
	for current != nil && current.createdAt.Before(time.Now().Add(-l.ttl)) {
		l.head = l.head.next
		l.currentLength--
		current = l.head
	}
}
