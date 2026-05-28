// Package audit implements audit event publishing
// and observer management.
package audit

// Observer processes audit events.
type Observer interface {
	Process(event Event) error
}
