package events

// Events are passed from modules to the decerver event handler. They should implement the Event
// interface. The event system is pub/sub. If you want an object to subscribe to events, make sure
// it implements the Subscriber interface and pass it to the event system.
import (
	"time"
)

type EventChan chan *Event

// This interface allow modules to subscribe to and publish events.
type EventProcessor interface {
	Post(e Event) error
	Subscribe(sub Subscriber) error
	Unsubscribe(id string) error
	TrafficData() string
}

// A default object that implements 'Event'
type Event struct {
	Event     string
	Target    string
	Resource  interface{}
	Source    string
	TimeStamp time.Time
}

// A subscriber listens to events.
type Subscriber interface {
	// Post an event
	Post(Event)
	// The subscriber only listen to events published by this source
	Source() string
	// The subscriber Id (must be unique).
	Id() string
	// The type of event it subscribes to.
	Event() string
	// The target (if any).
	Target() string
}
