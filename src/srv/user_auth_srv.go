package srv

import (
	"errors"
	"fmt"
	"time"

	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/golang-jwt/jwt"
)

type UserAuth interface {
	Jwt() string
}

type UserAuthSrv struct {
	secretKey string
	repo      repo.UserAuthRepo
}

func UserAuthSrvCtor(secretKey string, repo repo.UserAuthRepo) UserAuthSrv {
	return UserAuthSrv{secretKey, repo}
}

func (u UserAuthSrv) Jwt(Username, Password string) (string, error) {
	userId, err := u.repo.UserId(Username)
	passwordHash, err := u.repo.PasswordHash(Username)
	passValid := CheckPasswordHash(Password, passwordHash)
	if !passValid {
		return "", errors.New("Invalid password")
	}
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.secretKey))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error generate jwt token: '%s'", err))
	}
	return t, nil
}
