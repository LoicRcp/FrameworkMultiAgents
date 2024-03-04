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
	RegisterContainerAnswerContent
	RegisterAgentContent
	InterAgentMessageContent
)

type Message struct {
	Type           MessageType
	Sender         string
	ContentType    ContentType
	Content        string // Serialized content
	CorrelationID  int64  // Unique ID for matching requests and responses
	ExpectResponse bool   `json:"expectResponse"`
}

type RegisterContainerPayload struct {
	Address string
}

type RegisterContainerAnswerPayload struct {
	ID    string
	Error string
}

type RegisterAgentPayload struct {
	ContainerID string
}

type InterAgentMessagePayload struct {
	ReceiverID int
	Content    string
}

func (registerContainerPayload RegisterContainerPayload) String() string {
	return registerContainerPayload.Address
}

func (registerContainerAnswerPayload RegisterContainerAnswerPayload) String() string {
	return registerContainerAnswerPayload.Error
}

func (registerAgentPayload RegisterAgentPayload) String() string {
	return registerAgentPayload.ContainerID
}
