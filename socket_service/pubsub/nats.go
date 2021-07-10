package pubsub

import (
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/nats-io/nats.go"
)

type NATSPubSub struct {
	conn          *nats.EncodedConn
	subscriptions map[string]*nats.Subscription
}

var _ service.IPubSub = &NATSPubSub{}

func NewNatsPubSub(config service.Config) *NATSPubSub {
	rawConn, err := nats.Connect(config.PubSubConn)
	if err != nil {
		log.Fatalf("fatal: error connecting to NATS: %v", err)
	}
	conn, err := nats.NewEncodedConn(rawConn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("fatal: error setting NATS encoder: %v", err)
	}
	return &NATSPubSub{
		conn,
		make(map[string]*nats.Subscription),
	}
}

func (n *NATSPubSub) PublishMessage(msg service.Message) {
	err := n.conn.Publish(msg.Room, msg)
	if err != nil {
		log.Printf("error: publishing message: %v, got error: %v", msg, err)
	}
}

func (n *NATSPubSub) StartListening(room string, callback func(msg service.Message)) {
	sub, err := n.conn.Subscribe(room, func(msg *service.Message) {
		callback(*msg)
	})
	if err != nil {
		log.Printf("error: subscribing to topic %v: %v", room, err)
		return
	}
	n.subscriptions[room] = sub
}

func (n *NATSPubSub) StopListening(room string) {
	sub, ok := n.subscriptions[room]
	if !ok {
		log.Println("error: tried to unsubscribe for a non existing topic")
		return
	}
	err := sub.Unsubscribe()
	if err != nil {
		log.Printf("error: unsubscribing from topic %v: %v", room, err)
	}
	delete(n.subscriptions, room)
}
