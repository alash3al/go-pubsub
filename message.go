package pubsub

// Message The message metadata
type Message struct {
	topic     string
	payload   interface{}
	createdAt int64
}

// GetTopic Return the topic of the current message
func (m *Message) GetTopic() string {
	return m.topic
}

// GetPayload Get the payload of the current message
func (m *Message) GetPayload() interface{} {
	return m.payload
}

// GetCreatedAt Get the creation time of this message
func (m *Message) GetCreatedAt() int64 {
	return m.createdAt
}
