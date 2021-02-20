package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// Reference: https://gist.githubusercontent.com/eamonnmcevoy/c7ab5a5253712561f8dd923936646b96/raw/608ae87baf3a68053e7724dc0ab1bf10789587d4/hash.go

type PasswordHandler interface {
	//Generate a salted hash for the input string
	Generate(s string) (string, error)
	//Compare string to generated hash
	Compare(hash string, s string) error
}

type Hash struct{}

func (c *Hash) Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func (c *Hash) Compare(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
