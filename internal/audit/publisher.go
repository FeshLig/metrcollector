package audit

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
