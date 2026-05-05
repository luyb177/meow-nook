package jwtx

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

const (
	RoleUser      = "user"
	RoleVolunteer = "volunteer"
	RoleAdmin     = "admin"
)

type Claims struct {
	UserID int64
	Role   string
}

func Parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	mc, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	var userID int64
	switch v := mc["user_id"].(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	default:
		return nil, errors.New("invalid user_id in token")
	}

	role, _ := mc["role"].(string)

	return &Claims{UserID: userID, Role: role}, nil
}
