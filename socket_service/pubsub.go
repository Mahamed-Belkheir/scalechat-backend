package socket_service

type IPubSub interface {
	startListening(room string, handler func(msg Message))
	stopListening(room string)
	publishMessage(msg Message)
}
