package messenger

import (
	"fmt"
	"gotil/loadbalancer/message"
	"io"
	"os"
)

type Sender interface {
	Send(msg message.Message)
	SetLogger(writer io.Writer)
}

type messenger struct {
	id       string
	outgoing chan<- message.Message
	log      io.Writer
}

func NewSender(id string, outgoing chan<- message.Message) Sender {
	return &messenger{
		id:       id,
		outgoing: outgoing,
		log:      os.Stdout,
	}
}

func (m *messenger) Send(msg message.Message) {
	m.outgoing <- msg
	fmt.Fprintf(m.log, "%s sent %v", m.id, msg)
}

func (m *messenger) SetLogger(writer io.Writer) {
	m.log = writer
}
