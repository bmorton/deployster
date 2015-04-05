package poller

type Handler interface {
	Handle(*Event)
}

type HandlerFunc func(*Event)

func (f HandlerFunc) Handle(e *Event) {
	f(e)
}
