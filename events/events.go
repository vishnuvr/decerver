package events

import (
	"fmt"
	"github.com/eris-ltd/decerver-interfaces/core"
	"github.com/eris-ltd/decerver-interfaces/events"
	"strings"
	"sync"
)

type subscriptions struct {
	srs []events.Subscriber
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
	channels map[string]*subscriptions
	glob     *subscriptions
	ate      core.RuntimeManager
}

func NewEventProcessor(ate core.RuntimeManager) *EventProcessor {
	ep := &EventProcessor{}
	ep.channels = make(map[string]*subscriptions)
	ep.glob = NewSubscriptions()
	ep.ate = ate
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
	fmt.Println("Posting stuff " + e.Target)
	src := e.Source
	fmt.Println(src)

	subs := ep.glob
	for _, sub := range subs.srs {
		fmt.Println("Found service")
		fmt.Printf("Chan: %v\n", sub)
		sub.Channel() <- e
	}

	if src != "*" {
		subs = ep.channels[src]
		if subs == nil {
			ep.mutex.Unlock()
			return
		}
		for _, sub := range subs.srs {
			sub.Channel() <- e
		}
	}

	ep.mutex.Unlock()
}

/*
func (ep *EventProcessor) SubscribeNoChan(source string, callback string) {

}

	src := sub.Source()
	var split []string

	if strings.Trim(sub.Source(), " ") == "*" {
		ep.glob.srs = append(ep.glob.srs, sub)
		fmt.Println("Subscriber added to globals")
		fmt.Printf("SUUUUUBBBBBB %v\n",sub)
		return
	}

	split = strings.Split(src, ";")

	for _, s := range split {
		subs := ep.channels[s]
		if subs == nil {
			newSubs := NewSubscriptions()
			ep.channels[s] = newSubs
			subs = newSubs
		}
		subs.srs = append(subs.srs, sub)
		fmt.Printf("New subscriber added to: %s\n", sub.Source())
	}
}
*/
func (ep *EventProcessor) Subscribe(sub events.Subscriber) {

	src := sub.Source()
	var split []string

	if strings.Trim(sub.Source(), " ") == "*" {
		ep.glob.srs = append(ep.glob.srs, sub)
		fmt.Println("Subscriber added to globals")
		return
	}

	split = strings.Split(src, ";")

	for _, s := range split {
		subs := ep.channels[s]
		if subs == nil {
			newSubs := NewSubscriptions()
			ep.channels[s] = newSubs
			subs = newSubs
		}
		subs.srs = append(subs.srs, sub)
		fmt.Printf("New subscriber added to: %s\n", sub.Source())
	}
}

func (ep *EventProcessor) Unsubscribe(sub events.Subscriber) {

	src := sub.Source()

	var split []string
	id := sub.Id()
	if strings.Trim(sub.Source(), " ") == "*" {
		theIdx := -1
		for i, sub := range ep.glob.srs {
			if sub.Id() == id {
				theIdx = i
				break
			}
			if theIdx >= 0 {
				sub.Close()
				fmt.Println("Subscriber removed from globals")
				ep.glob.srs = append(ep.glob.srs[:theIdx], ep.glob.srs[theIdx+1:]...)
			}
		}
		return
	}

	split = strings.Split(src, ";")

	for _, s := range split {
		subs := ep.channels[s]
		theIdx := -1
		for i, sub := range subs.srs {
			if sub.Id() == id {
				theIdx = i
				break
			}
			if theIdx >= 0 {
				fmt.Printf("Subscriber removed from: %s, id: %s\n", sub.Source(), sub.Id())
				subs.srs = append(subs.srs[:theIdx], subs.srs[theIdx+1:]...)
			}
		}

	}

}
