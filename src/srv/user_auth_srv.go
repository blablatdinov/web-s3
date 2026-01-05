package srv

import (
	"errors"
	"fmt"
	"time"

	"github.com/blablatdinov/web-s3/src/repo"
	"github.com/golang-jwt/jwt"
)

type UserAuth interface {
	Jwt(string, string) (string, error)
	Validate(token string) (bool, error)
	ExtractClaims(token string) (jwt.MapClaims, error)
}

type UserAuthSrv struct {
	secretKey string
	repo      repo.UserAuthRepo
}

func UserAuthSrvCtor(secretKey string, repo repo.UserAuthRepo) UserAuth {
	return UserAuthSrv{secretKey, repo}
}

func (u UserAuthSrv) Jwt(Username, Password string) (string, error) {
	userId, err := u.repo.UserId(Username)
	if err != nil {
		return "", err
	}
	passwordHash, err := u.repo.PasswordHash(Username)
	if err != nil {
		return "", err
	}
	passValid := PswrdCtor(Password).Check(passwordHash)
	if !passValid {
		return "", errors.New("invalid password")
	}
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.secretKey))
	if err != nil {
		return "", fmt.Errorf("error generate jwt token: %w", err)
	}
	return t, nil
}

func (u UserAuthSrv) Validate(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.secretKey), nil
	})
	if err != nil {
		return false, err
	}
	if !parsedToken.Valid {
		return false, errors.New("token is not valid")
	}
	return true, nil
}

func (u UserAuthSrv) ExtractClaims(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
