package messaging

import (
	"github.com/google/uuid"
)

const TopicVerifyEmail = "notify.verify.email"

type VerifyEmail struct {
	UserUUID        uuid.UUID `msgpack:"id"`
	FirstName       string    `msgpack:"fn"`
	Email           string    `msgpack:"el"`
	Code            string    `msgpack:"ce"`
	VerificationURL string    `msgpack:"url"`
}
