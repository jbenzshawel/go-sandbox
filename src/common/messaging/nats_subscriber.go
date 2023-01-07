package messaging

import "github.com/nats-io/nats.go"

type NatsSubscriber struct {
	nc       *nats.Conn
	handlers map[string]func(data []byte)
}

func NewNatsSubscriber(nc *nats.Conn) *NatsSubscriber {
	return &NatsSubscriber{nc: nc, handlers: map[string]func(data []byte){}}
}

func (s *NatsSubscriber) WithSubscription(topic string, handler func(data []byte)) *NatsSubscriber {
	s.handlers[topic] = handler

	return s
}

func (s *NatsSubscriber) Subscribe() error {
	for topic, handler := range s.handlers {
		_, err := s.nc.Subscribe(topic, func(msg *nats.Msg) {
			handler(msg.Data)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
