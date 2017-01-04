package server

type Message struct {
	body string
}

func NewMessage(body string) *Message {
	return &Message{
		body,
	}
}
