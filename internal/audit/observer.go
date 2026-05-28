package audit

// Observer processes audit events.
type Observer interface {
	Process(event Event) error
}
