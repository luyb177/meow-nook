// Package jwtx provides shared JWT parsing and claims helpers used by the gateway.
package jwtx

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Role constants — must match the values stored in the users table.
const (
	RoleUser      = "user"
	RoleVolunteer = "volunteer"
	RoleAdmin     = "admin"
)

// Claims is the parsed representation of a meow-nook JWT.
type Claims struct {
	UserID int64
	Role   string
}

// Parse validates the token string with the given HMAC secret and returns Claims.
// Returns an error if the token is invalid, expired, or malformed.
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

	// user_id may be float64 (JSON number) or int64 depending on the JWT library version.
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
