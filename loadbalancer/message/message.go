package message

type messageType string

const NewTask = messageType("new_task")
const TaskCompleted = messageType("completed[new_task]")
const TaskStarted = messageType("started[new_taks]")
const WorkerDied = messageType("worker died")
const WorkerStarted = messageType("worker started")

type Message struct {
	Type    messageType
	Payload interface{}
}

func OfType(msgType messageType) Message {
	return Message{
		Type: msgType,
	}
}

func (m Message) WithPayload(payload interface{}) Message {
	m.Payload = payload
	return m
}
