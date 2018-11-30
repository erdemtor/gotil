package messenger

import (
	"gotil/loadbalancer/message"
	"log"
)

type Messenger interface {
	Send(msg message.Message)
	Receive() <-chan message.Message
}

type messenger struct {
	id       string
	outgoing chan<- message.Message
	incoming <-chan message.Message
}

func New(id string, incoming <-chan message.Message, outgoing chan<- message.Message) Messenger {
	m := &messenger{
		id:       id,
		outgoing: outgoing,
		incoming: incoming,
	}
	return m
}

func (m *messenger) Send(msg message.Message) {
	m.outgoing <- msg
	log.Printf("%s sent %s\n", m.id, msg)

}

func (m *messenger) Receive() <-chan message.Message {
	retChan := make(chan message.Message, 1)
	go func() {
		msg := <-m.incoming
		log.Printf("%s received %s \n", m.id, msg)
		retChan <- msg
		close(retChan)
	}()
	return retChan
}
