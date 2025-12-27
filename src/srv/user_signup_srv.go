package srv

import (
	"errors"
	"fmt"

	"github.com/blablatdinov/web-s3/src/repo"
)

var (
	ErrUsernameEmpty = errors.New("empty username")
)

type UserSignupSrv interface {
	Create(username, rawPassword string) (int, error)
}

type UsrSignupSrv struct {
	repo repo.UserSignupRepo
}

func UsrSignupSrvCtor(repo repo.UserSignupRepo) UserSignupSrv {
	return UsrSignupSrv{repo}
}

func (u UsrSignupSrv) Create(username, rawPassword string) (int, error) {
	if username == "" {
		return 0, fmt.Errorf("%w", ErrUsernameEmpty)
	}
	hashedPassword, err := PswrdCtor(rawPassword).Hash()
	if err != nil {
		return 0, err
	}
	userId, err := u.repo.Create(username, hashedPassword)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
