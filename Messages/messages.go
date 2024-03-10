package Messages

import "strconv"

type MessageType int
type ContentType int

const (
	RegisterContainer MessageType = iota
	RegisterContainerAnswer
	RegisterAgent
	RegisterAgentAnswer
	InterAgentAsyncMessage
	GetAgentAdress
	GetAgentAdressAnswer
)

const (
	RegisterContainerContent ContentType = iota
	RegisterContainerAnswerContent
	RegisterAgentContent
	RegisterAgentAnswerContent
	InterAgentAsyncMessageContent
	GetAgentAdressContent
	GetAgentAdressAnswerContent
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
	ContainerID string
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

type GetAgentAdressPayload struct {
	AgentID string
}

type GetAgentAdressAnswerPayload struct {
	Adress string
}

func (registerContainerPayload RegisterContainerPayload) String() string {
	return registerContainerPayload.Address
}

func (registerContainerAnswerPayload RegisterContainerAnswerPayload) String() string {
	return registerContainerAnswerPayload.ContainerID
}

func (registerAgentPayload RegisterAgentPayload) String() string {
	return registerAgentPayload.ContainerID
}

func (interAgentAsyncMessagePayload InterAgentAsyncMessagePayload) String() string {
	return interAgentAsyncMessagePayload.Content
}

func (registerAgentAnswerPayload RegisterAgentAnswerPayload) String() string {
	return strconv.Itoa(registerAgentAnswerPayload.ID)
}

func (getAgentAdressPayload GetAgentAdressPayload) String() string {
	return getAgentAdressPayload.AgentID
}

func (getAgentAdressAnswerPayload GetAgentAdressAnswerPayload) String() string {
	return getAgentAdressAnswerPayload.Adress
}

func (message Message) String() string {
	return message.Sender + message.Content
}
