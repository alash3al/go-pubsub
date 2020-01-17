package pubsub

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

// Subscribers ...
type Subscribers map[string]*Subscriber

// Subscriber ...
type Subscriber struct {
	id        string
	messages  chan *Message
	createdAt int64
	destroyed bool
	topics    map[string]bool
	lock      sync.RWMutex
}

// NewSubscriber returns a new subscriber
func NewSubscriber() (*Subscriber, error) {
	id := make([]byte, 50)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}

	return &Subscriber{
		id:        hex.EncodeToString(id),
		messages:  make(chan *Message),
		createdAt: time.Now().UnixNano(),
		destroyed: false,
		lock:      sync.RWMutex{},
		topics:    map[string]bool{},
	}, nil
}

// GetID return the subscriber id
func (s *Subscriber) GetID() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.id
}

// GetCreatedAt return `time.Time` of the creation time
func (s *Subscriber) GetCreatedAt() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.createdAt
}

// AddTopic Adds a new topic to the subscriber
func (s *Subscriber) AddTopic(topic string) {
	s.lock.Lock()
	s.topics[topic] = true
	s.lock.Unlock()
}

// RemoveTopic Removes a topic from the subscriber
func (s *Subscriber) RemoveTopic(topic string) {
	s.lock.Lock()
	delete(s.topics, topic)
	s.lock.Unlock()
}

// GetTopics return slice of subscriber topics
func (s *Subscriber) GetTopics() []string {
	s.lock.RLock()
	subscriberTopics := s.topics
	s.lock.RUnlock()

	topics := []string{}
	for topic := range subscriberTopics {
		topics = append(topics, topic)
	}
	return topics
}

// GetMessages returns a channel of *Message to listen on
func (s *Subscriber) GetMessages() <-chan *Message {
	return s.messages
}

// Signal sends a message to subscriber
func (s *Subscriber) Signal(m *Message) *Subscriber {
	s.lock.RLock()
	if !s.destroyed {
		s.messages <- m
	}
	s.lock.RUnlock()

	return s
}

// close the underlying channels/resources
func (s *Subscriber) destroy() {
	s.lock.Lock()
	s.destroyed = true
	s.lock.Unlock()

	close(s.messages)
}
