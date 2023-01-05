package messaging

import "github.com/google/uuid"

const TOPIC_VERIFY_EMAIL = "notify.verify.email"

type VerifyEmail struct {
	UserUUID        uuid.UUID
	Email           string
	Code            string
	VerificationURL string
}
