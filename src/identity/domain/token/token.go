package token

import (
	"crypto/rand"

	"github.com/google/uuid"
)

type Token struct {
	code string
}

func NewToken() (*Token, error) {
	code, err := generateCode()
	if err != nil {
		return nil, err
	}
	return &Token{
		code: code,
	}, nil
}

const tokenChars = "1234567890"

func generateCode() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	tokenCharsLength := len(tokenChars)
	for i := 0; i < 6; i++ {
		buffer[i] = tokenChars[int(buffer[i])%tokenCharsLength]
	}

	return string(buffer), nil
}

func (t *Token) Code() string {
	return t.code
}

func Verify(repo Repository, key uuid.UUID, code string) bool {
	token := repo.GetToken(key)
	if token == nil {
		return false
	}

	isValid := code == token.Code()
	if isValid {
		repo.ClearToken(key)
	}
	return isValid
}
