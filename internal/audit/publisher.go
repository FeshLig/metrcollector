package audit

// Closer closes observer resources.
type Closer interface {
	Close() error
}

// Publisher distributes audit events to observers.
type Publisher struct {
	observers []Observer
}

// NewPublisher creates new audit publisher.
func NewPublisher() *Publisher {
	return &Publisher{}
}

// Subscribe registers new audit observer.
func (p *Publisher) Subscribe(observer Observer) {
	p.observers = append(p.observers, observer)
}

// Notify sends audit event to all observers.
func (p *Publisher) Notify(event Event) {
	for _, observer := range p.observers {
		_ = observer.Process(event)
	}
}

// Close closes all observers that implement Closer interface.
func (p *Publisher) Close() error {
	for _, observer := range p.observers {

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
