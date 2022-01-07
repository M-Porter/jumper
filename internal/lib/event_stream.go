package lib

type Event = int

const (
	EventUpdate Event = iota
	EventDone
)

type Events chan Event

func (e Events) Close() {
	close(e)
}

func (e Events) Update() {
	e <- EventUpdate
}

func (e Events) Done() {
	e <- EventDone
}

var sharedEvent Events

func EventsStream() Events {
	return sharedEvent
}

func init() {
	sharedEvent = make(Events)
}
