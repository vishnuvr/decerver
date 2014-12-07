package events

import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/events"
	"github.com/eris-ltd/decerver-interfaces/modules"
	"sync"
)

type SubMap map[string]*subscriptions

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
		// DEBUG
		fmt.Println("Subscriber removed: " + subId)
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
	mutex          *sync.Mutex
	postChan       chan events.Event
	channels       map[string]SubMap
	byId           map[string]events.Subscriber
	moduleRegistry modules.ModuleRegistry
}

func NewEventProcessor(mr modules.ModuleRegistry) *EventProcessor {
	ep := &EventProcessor{}
	ep.mutex = &sync.Mutex{}
	ep.channels = make(map[string]SubMap)
	// DEBUG
	fmt.Printf("[Events] Subscriber map created: %v\n", ep.channels)
	ep.byId = make(map[string]events.Subscriber)
	ep.postChan = make(chan events.Event, 64)
	ep.moduleRegistry = mr
	return ep
}

func (ep *EventProcessor) Post(e events.Event) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	src := e.Source
	// DEBUG
	fmt.Println("Receiving event '" + e.Event + "' from '" + e.Source + "'.")

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
		if sub.Target() == e.Target {
			fmt.Println("Found service")
			fmt.Printf("Chan: %v\n", sub)
			sub.Channel() <- e
		}
	}
}

func (ep *EventProcessor) Subscribe(sub events.Subscriber) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()
	src := sub.Source()
	fmt.Println("[Events] New subscription registering: " + src)
	srcSubs, okSrc := ep.channels[src]
	if !okSrc {
		srcSubs = make(SubMap)
		ep.channels[src] = srcSubs
	}

	evt := sub.Event()
	evts, okEvt := srcSubs[evt]

	if !okEvt {
		evts = NewSubscriptions()
		srcSubs[evt] = evts
	}

	evts.add(sub)
	ep.byId[sub.Id()] = sub

	sub.SetChannel(ep.moduleRegistry.GetModules()[src].Subscribe(sub.Id(), sub.Event(), sub.Target()))
	fmt.Printf("New subscriber added to: %s (%s)\n", sub.Source(), sub.Event())
}

func (ep *EventProcessor) Unsubscribe(id string) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()
	sub, ok := ep.byId[id]
	if !ok {
		fmt.Println("No subscriber with id: " + id)
		return
	}
	ep.moduleRegistry.GetModules()[sub.Source()].UnSubscribe(sub.Id())
	ep.channels[sub.Source()][sub.Event()].remove(id)
	ep.byId[id] = nil
}
