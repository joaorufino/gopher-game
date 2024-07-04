package event

import (
	"container/heap"
	"sync"

	"github.com/joaorufino/cv-game/internal/interfaces"
)

// EventManager manages event registration and dispatching with priority and async handling.
type EventManager struct {
	handlers   map[interfaces.EventType][]interfaces.EventHandler
	mu         sync.RWMutex
	eventQueue PriorityQueue
	wg         sync.WaitGroup
}

// NewEventManager creates a new event dispatcher.
func NewEventManager() *EventManager {
	return &EventManager{
		handlers:   make(map[interfaces.EventType][]interfaces.EventHandler),
		eventQueue: make(PriorityQueue, 0),
	}
}

// RegisterHandler registers an event handler for a specific event type.
func (d *EventManager) RegisterHandler(eventType interfaces.EventType, handler interfaces.EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch dispatches an event to all registered handlers with priority.
func (d *EventManager) Dispatch(event interfaces.Event) {
	d.mu.Lock()
	heap.Push(&d.eventQueue, &Item{
		value:    event,
		priority: event.Priority,
	})
	d.mu.Unlock()
	d.wg.Add(1)
	go d.processEvents()
}

// processEvents processes events from the priority queue asynchronously.
func (d *EventManager) processEvents() {
	defer d.wg.Done()

	var event interfaces.Event

	d.mu.Lock()
	if d.eventQueue.Len() > 0 {
		item := heap.Pop(&d.eventQueue).(*Item)
		event = item.value.(interfaces.Event)
	}
	d.mu.Unlock()

	if event.Type != "" {
		d.mu.RLock()
		if handlers, exists := d.handlers[event.Type]; exists {
			for _, handler := range handlers {
				handler(event)
			}
		}
		d.mu.RUnlock()
	}
}

// Wait waits for all events to be processed.
func (d *EventManager) Wait() {
	d.wg.Wait()
}

// PriorityQueue implements a priority queue for events.
type PriorityQueue []*Item

// Item represents an item in the priority queue.
type Item struct {
	value    interface{}
	priority int
	index    int
}

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares the priority of two items.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

// Swap swaps two items in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the highest priority item from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
