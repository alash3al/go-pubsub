package pubsub

type Message struct {
	topic     string
	payload   interface{}
	createdAt int64
}

func (m *Message) GetTopic() string {
	return m.topic
}

func (m *Message) GetPayload() interface{} {
	return m.payload
}

func (m *Message) GetCreatedAt() int64 {
	return m.createdAt
}
