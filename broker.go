package pubsub

import (
	"sync"
	"time"
)

// Broker The broker related meta data
type Broker struct {
	subscribers Subscribers
	sLock       sync.RWMutex

	topics map[string]Subscribers
	tLock  sync.RWMutex
}

// NewBroker Create new broker
func NewBroker() *Broker {
	return &Broker{
		subscribers: Subscribers{},
		sLock:       sync.RWMutex{},
		topics:      map[string]Subscribers{},
		tLock:       sync.RWMutex{},
	}
}

// Attach Create a new subscriber and register it into our main broker
func (b *Broker) Attach() (*Subscriber, error) {
	s, err := NewSubscriber()

	if err != nil {
		return nil, err
	}

	b.sLock.Lock()
	b.subscribers[s.GetID()] = s
	b.sLock.Unlock()

	return s, nil
}

// Subscribe subscribes the specified subscriber "s" to the specified list of topic(s)
func (b *Broker) Subscribe(s *Subscriber, topics ...string) {
	b.tLock.Lock()
	defer b.tLock.Unlock()

	for _, topic := range topics {
		if nil == b.topics[topic] {
			b.topics[topic] = Subscribers{}
		}
		s.AddTopic(topic)
		b.topics[topic][s.id] = s
	}

}

// Unsubscribe Unsubscribe the specified subscriber from the specified topic(s)
func (b *Broker) Unsubscribe(s *Subscriber, topics ...string) {
	for _, topic := range topics {
		b.tLock.Lock()
		if nil == b.topics[topic] {
			continue
		}
		delete(b.topics[topic], s.id)
		b.tLock.Unlock()
		s.RemoveTopic(topic)
	}
}

// Detach remove the specified subscriber from the broker
func (b *Broker) Detach(s *Subscriber) {
	s.destroy()
	b.sLock.Lock()
	b.Unsubscribe(s, s.GetTopics()...)
	delete(b.subscribers, s.id)
	defer b.sLock.Unlock()
}

// Broadcast broadcast the specified payload to all the topic(s) subscribers
func (b *Broker) Broadcast(payload interface{}, topics ...string) {
	for _, topic := range topics {
		if b.Subscribers(topic) < 1 {
			continue
		}
		b.tLock.RLock()
		for _, s := range b.topics[topic] {
			m := &Message{
				topic:     topic,
				payload:   payload,
				createdAt: time.Now().UnixNano(),
			}
			go (func(s *Subscriber) {
				s.Signal(m)
			})(s)
		}
		b.tLock.RUnlock()
	}
}

// Subscribers Get the subscribers count
func (b *Broker) Subscribers(topic string) int {
	b.tLock.RLock()
	defer b.tLock.RUnlock()
	return len(b.topics[topic])
}

// GetTopics Returns a slice of topics
func (b *Broker) GetTopics() []string {
	b.tLock.RLock()
	brokerTopics := b.topics
	b.tLock.RUnlock()

	topics := []string{}
	for topic := range brokerTopics {
		topics = append(topics, topic)
	}

	return topics
}
