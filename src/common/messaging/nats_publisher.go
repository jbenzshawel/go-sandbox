package messaging

import "github.com/nats-io/nats.go"

type Publisher interface {
	Publish(topic string, msg []byte) error
}

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) *NatsPublisher {
	return &NatsPublisher{nc: nc}
}

func (p *NatsPublisher) Publish(topic string, msg []byte) error {
	return p.nc.Publish(topic, msg)
}
