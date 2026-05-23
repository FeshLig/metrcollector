package audit

import "sync"

const maxConcurrentObservers = 8

// Closer closes observer resources.
type Closer interface {
	Close() error
}

// Publisher distributes audit events to observers.
type Publisher struct {
	observers []Observer
	mu        sync.RWMutex
	sem       chan struct{}
}

// NewPublisher creates new audit publisher.
func NewPublisher() *Publisher {
	return &Publisher{
		sem: make(chan struct{}, maxConcurrentObservers),
	}
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

	for _, observer := range observers {
		wg.Add(1)

		p.sem <- struct{}{}

		go func(observer Observer) {
			defer wg.Done()
			defer func() { <-p.sem }()

			_ = observer.Process(event)
		}(observer)
	}

	wg.Wait()
}

// Close closes all observers that implement Closer interface.
func (p *Publisher) Close() error {
	p.mu.RLock()

	observers := make([]Observer, len(p.observers))
	copy(observers, p.observers)

	p.mu.RUnlock()

	for _, observer := range observers {
		closer, ok := observer.(Closer)
		if !ok {
			continue
		}

		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}
