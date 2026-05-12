package audit

type Closer interface {
	Close() error
}

type Publisher struct {
	observers []Observer
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Subscribe(observer Observer) {
	p.observers = append(p.observers, observer)
}

func (p *Publisher) Notify(event Event) {
	for _, observer := range p.observers {
		_ = observer.Process(event)
	}
}

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
