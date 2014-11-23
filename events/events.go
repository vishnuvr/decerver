package events

import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/events"
	"sync"
)

type subscriptions struct {
	srs []events.Subscriber
}

func (ss *subscriptions) add(sub events.Subscriber) {
	ss.srs = append(ss.srs, sub)
}

func (ss *subscriptions) remove(subId string) bool {
	theIdx := -1
	for i, sub := range ss.srs {
		if sub.Id() == subId {
			theIdx = i
			break
		}
	}
	if theIdx >= 0 {
		sub := ss.srs[theIdx]
		sub.Close()
		// DEBUG
		fmt.Println("Subscriber removed from globals")
		ss.srs = append(ss.srs[:theIdx], ss.srs[theIdx+1:]...)
		return true
	}
	return false
}

func NewSubscriptions() *subscriptions {
	ss := &subscriptions{}
	ss.srs = make([]events.Subscriber, 0)
	return ss
}

type EventProcessor struct {
	mutex    sync.Mutex
	postChan chan events.Event
	// Sorts by source, then by event name.
	channels map[string]map[string]*subscriptions
	byId map[string]events.Subscriber
}

func NewEventProcessor() *EventProcessor {
	ep := &EventProcessor{}
	ep.channels = make(map[string]map[string]*subscriptions)
	ep.byId = make(map[string]events.Subscriber)
	ep.postChan = make(chan events.Event)
	return ep
}

func (ep *EventProcessor) Post(e events.Event) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	src := e.Source

	fmt.Println("Posting stuff " + e.Target)

	sourceSubs := ep.channels[src]
	if sourceSubs == nil {
		return
	}

	tgt := e.Target
	targetSubs := sourceSubs[tgt]
	if targetSubs == nil {
		return
	}

	for _, sub := range targetSubs.srs {
		// TODO necessary?
		if sub.Target() == e.Target {
			fmt.Println("Found service")
			fmt.Printf("Chan: %v\n", sub)
			sub.Channel() <- e
		}
	}
}

func (ep *EventProcessor) Subscribe(sub events.Subscriber) {
	src := sub.Source()
	srcSubs := ep.channels[src]
	if srcSubs == nil {
		srcSubs := make(map[string]*subscriptions)
		ep.channels[src] = srcSubs
		
	}
	evt := sub.Event()
	evts := srcSubs[evt]
	
	if evts == nil {
		evts := NewSubscriptions()
		srcSubs[evt] = evts
	}
	evts.add(sub)
	ep.byId[sub.Id()] = sub
	fmt.Printf("New subscriber added to: %s (%s)\n", sub.Source(), sub.Event())
}

func (ep *EventProcessor) Unsubscribe(id string) {
	sub, ok := ep.byId[id]
	if !ok {
		fmt.Println("No subscriber with id: " + id)
		return
	}
	ep.channels[sub.Source()][sub.Event()].remove(id)
	ep.byId[id] = nil
}
