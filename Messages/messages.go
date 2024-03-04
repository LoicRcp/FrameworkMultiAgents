package Messages

type MessageType int
type ContentType int

const (
	RegisterContainer MessageType = iota
	RegisterContainerAnswer
	RegisterAgent
	RegisterAgentAnswer
	InterAgentAsyncMessage
)

const (
	RegisterContainerContent ContentType = iota
	RegisterContainerAnswerContent
	RegisterAgentContent
	RegisterAgentAnswerContent
	InterAgentAsyncMessageContent
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
	Error string
}

type RegisterAgentPayload struct {
	ContainerID string
}

type InterAgentAsyncMessagePayload struct {
	ReceiverID int
	Content    string
}

type RegisterAgentAnswerPayload struct {
	ID int
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

func (interAgentAsyncMessagePayload InterAgentAsyncMessagePayload) String() string {
	return interAgentAsyncMessagePayload.Content
}

func (registerAgentAnswerPayload RegisterAgentAnswerPayload) String() string {
	return string(registerAgentAnswerPayload.ID)
}

func (message Message) String() string {
	return message.Sender + message.Content
}
