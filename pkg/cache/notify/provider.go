package notify

const notifyTopic = "cache:change:notify"

type IProvider interface {
	OK() bool
	Pub(channel string, payload []string) error
	Sub(channel string, callback func(payload []string)) error
}
