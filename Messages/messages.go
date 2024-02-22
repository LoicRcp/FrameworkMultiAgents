package Messages

type MessageType int
type ContentType int

const (
	RegisterContainer MessageType = iota
	RegisterContainerAnswer
	RegisterAgent
	InterAgentMessage
)

const (
	RegisterContainerContent ContentType = iota
)

type Message struct {
	Type          MessageType
	Sender        string
	ContentType   ContentType
	Content       string // Serialized content
	CorrelationID int64  // Unique ID for matching requests and responses
}

type RegisterContainerPayload struct {
	Address string
}

type RegisterContainerAnswerPayload struct {
	ID    string
	Error string
}

func (registerContainerPayload RegisterContainerPayload) String() string {
	return registerContainerPayload.Address
}

func (registerContainerAnswerPayload RegisterContainerAnswerPayload) String() string {
	return registerContainerAnswerPayload.ID + " " + registerContainerAnswerPayload.Error
}
