package events

import (
	"github.com/eris-ltd/deCerver-interfaces/events"
	"sync"
	"strings"
)

type subscriptions struct {
	srs []events.Subscriber
}

func NewSubscriptions() *subscriptions {
	ss := &subscriptions{}
	ss.srs = make([]events.Subscriber, 1)
	return ss
}

type EventProcessor struct {
	mutex    sync.Mutex
	postChan chan events.Event
	// Sorts by source, then by event name.
	channels map[string]*subscriptions
	glob *subscriptions
}

func NewEventProcessor() *EventProcessor {
	ep := &EventProcessor{}
	ep.channels = make(map[string]*subscriptions)
	ep.glob = NewSubscriptions()
	/*
	for _ , mod := range modules {
		ep.channels[mod] = NewSubscriptions()
	} 
	*/
	
	ep.postChan = make(chan events.Event)
	return ep
}

func (ep *EventProcessor) Post(e events.Event) {
	ep.mutex.Lock()

	src := e.Source()
	subs := ep.channels[src]
	if subs == nil {
		return
	}
	for _, sub := range subs.srs {
		sub.Channel() <- e
	}
	
	subs = ep.glob
	for _, sub := range subs.srs {
		sub.Channel() <- e
	}
	
	ep.mutex.Unlock()
}

func (ep *EventProcessor) Subscribe(sub events.Subscriber) {
	src := sub.Source()
	var split []string
	
	if strings.Trim(sub.Source()," ") == "*" {
		ep.glob.srs = append(ep.glob.srs, sub)
	}
	
	split = strings.Split(src,";")

	for _ , s := range split {
		subs := ep.channels[s]
		subs.srs = append(subs.srs, sub)
	}

	
}

func (ep *EventProcessor) Unsubscribe(sub events.Subscriber) {
	
	src := sub.Source()
	
	var split []string
	ch := sub.Channel()
	if strings.Trim(sub.Source()," ") == "*" {
		theIdx := -1
		for i , sub := range ep.glob.srs {
			if sub.Channel() == ch {
				theIdx = i
				break
			}
			if theIdx >= 0 {
				ep.glob.srs = append(ep.glob.srs[:theIdx], ep.glob.srs[theIdx+1:]...)
			}
		}
		
	}
	
	split = strings.Split(src,";")
	
	for _ , s := range split {
		subs := ep.channels[s]
		theIdx := -1
		for i , sub := range subs.srs {
			if sub.Channel() == ch {
				theIdx = i
				break
			}
			if theIdx >= 0 {
				subs.srs = append(subs.srs[:theIdx], subs.srs[theIdx+1:]...)
			}
		}
		
	}

}
