package socket_service

type IPubSub interface {
	StartListening(room string, handler func(msg Message))
	StopListening(room string)
	PublishMessage(msg Message)
}
