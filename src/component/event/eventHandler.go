package event

type Handler interface {
	Handle(e Event) error
}
