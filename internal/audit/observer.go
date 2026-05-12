package audit

type Observer interface {
	Process(event Event) error
}
