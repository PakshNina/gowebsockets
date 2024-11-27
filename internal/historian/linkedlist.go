package historian

import (
	"sync"
	"time"
)

type node struct {
	value     interface{}
	createdAt time.Time
	next      *node
}

type linkedlist struct {
	head          *node
	tail          *node
	maxLength     int
	currentLength int
	ttl           time.Duration
	mutex         *sync.Mutex
}

func newLinkedList(maxLength int, ttl time.Duration) *linkedlist {
	return &linkedlist{
		maxLength: maxLength,
		ttl:       ttl,
		mutex:     &sync.Mutex{},
	}
}

func (l *linkedlist) saveValue(value interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	n := &node{
		value:     value,
		next:      nil,
		createdAt: time.Now(),
	}
	l.checkTTL()
	if l.head == nil {
		l.head = n
		l.tail = n
		l.currentLength++
		return
	}
	if l.currentLength+1 > l.maxLength {
		l.head = l.head.next
	}
	l.tail.next = n
	l.tail = l.tail.next
	l.currentLength++
}

func (l *linkedlist) getAllValues() []interface{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.checkTTL()
	current := l.head
	var result []interface{}
	for current != nil {
		result = append(result, current.value)
		current = current.next
	}
	return result
}

func (l *linkedlist) checkTTL() {
	current := l.head
	for current != nil && current.createdAt.Before(time.Now().Add(-l.ttl)) {
		l.head = l.head.next
		current = l.head
		l.currentLength--
	}
}
