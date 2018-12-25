package message

import "fmt"

type messageType string

// Predefined message type constants
const (
	NewTask       = messageType("new_task")
	TaskCompleted = messageType("completed[new_task]")
	TaskStarted   = messageType("started[new_task]")
	WorkerDied    = messageType("worker died")
	WorkerStarted = messageType("worker started")
)

//Message is a data struct to enable the communication between master & workers
type Message struct {
	Type    messageType
	Payload interface{}
}

func (m Message) String() string {
	return fmt.Sprintf("{type: %s}", m.Type)
}

//OfType creates a new message of given ype
func OfType(msgType messageType) Message {
	return Message{
		Type: msgType,
	}
}

//WithPayload adds a payload to the message that's constructed already and returns it.
func (m Message) WithPayload(payload interface{}) Message {
	m.Payload = payload
	return m
}
