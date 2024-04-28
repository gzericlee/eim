package mq

type Topic string

const (
	GroupMessageDispatchTopic Topic = "group_message_dispatch"
	UserMessageDispatchTopic  Topic = "user_message_dispatch"
	MessageSendDispatchTopic  Topic = "%s_send"
)

var MessageTopics = map[string]Topic{
	//ToUser
	"1": UserMessageDispatchTopic,

	//ToGroup
	"2": GroupMessageDispatchTopic,
}
