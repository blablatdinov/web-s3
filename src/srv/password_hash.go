package srv

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password interface {
	Hash() (string, error)
	Check(hash string) bool
}

var (
	ErrorHashingPassword = errors.New("error hashing password")
)

type Pswrd struct {
	rawPass string
}

func PswrdCtor(rawPass string) Password {
	return Pswrd{rawPass}
}

func (p Pswrd) Hash() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.rawPass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w %s", ErrorHashingPassword, err)
	}
	return string(hash), nil
}

func (p Pswrd) Check(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p.rawPass))
	return err == nil
}
