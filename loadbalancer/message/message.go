package message

type messageType string

const NewTask = messageType("new NewTask")
const TaskCompleted = messageType("NewTask completed")
const TaskStarted = messageType("NewTask started")
const WorkerDied = messageType("worker died")
const WorkerStarted = messageType("worker started")

type Message struct {
	Type messageType
	Data interface{}
}

func OfType(msgType messageType) Message {
	return Message{
		Type: msgType,
	}
}

func (m Message) WithData(data interface{}) Message {
	m.Data = data
	return m
}
