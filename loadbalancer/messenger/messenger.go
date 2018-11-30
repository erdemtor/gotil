package messenger

import (
	"gotil/loadbalancer/message"
	"log"
)

type Sender interface {
	Send(msg message.Message)
}

type messenger struct {
	id       string
	outgoing chan<- message.Message
	incoming <-chan message.Message
}

func NewSender(id string, outgoing chan<- message.Message) Sender {
	return &messenger{
		id:       id,
		outgoing: outgoing,
	}
}

func (m *messenger) Send(msg message.Message) {
	m.outgoing <- msg
	log.Printf("%s sent %v", m.id, msg)
}
