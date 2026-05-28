package audit

import "sync"

const maxConcurrentObservers = 8

// Publisher distributes audit events to observers.
type Publisher struct {
	observers []Observer
	mu        sync.RWMutex
}

// NewPublisher creates new audit publisher.
func NewPublisher() *Publisher {
	return &Publisher{}
}

// Subscribe registers new audit observer.
func (p *Publisher) Subscribe(observer Observer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.observers = append(p.observers, observer)
}

// Notify sends audit event to all observers.
func (p *Publisher) Notify(event Event) {
	p.mu.RLock()

	observers := make([]Observer, len(p.observers))
	copy(observers, p.observers)

	p.mu.RUnlock()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentObservers)

	for _, observer := range observers {
		wg.Add(1)

		sem <- struct{}{}

		go func(observer Observer) {
			defer wg.Done()
			defer func() { <-sem }()

			_ = observer.Process(event)
		}(observer)
	}

	wg.Wait()
}
