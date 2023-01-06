package publisher

import (
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/jbenzshawel/go-sandbox/common/messaging"
)

type NatsPublisher struct {
	natsURL string
}

func NewNatsPublisher(natsURL string) *NatsPublisher {
	return &NatsPublisher{natsURL: natsURL}
}

func (p *NatsPublisher) NotifyVerifyEmailPublisher() func(msg *messaging.VerifyEmail) error {
	return func(msg *messaging.VerifyEmail) error {
		nc, err := nats.Connect(p.natsURL)
		if err != nil {
			return err
		}

		msgBytes, err := msgpack.Marshal(msg)
		if err != nil {
			return err
		}

		return nc.Publish(messaging.TOPIC_VERIFY_EMAIL, msgBytes)
	}
}
