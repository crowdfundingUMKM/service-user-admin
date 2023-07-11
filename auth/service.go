package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userUnixID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

var secret = os.Getenv("JWT_SECRET")
var SECRET_KEY = []byte(secret)

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userUnixID string) (string, error) {
	claim := jwt.MapClaims{}
	claim["unix_id"] = userUnixID
	claim["exp"] = time.Now().Add(2 * 24 * time.Hour).Unix() // Expires in 2 days (48 hours)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil

}

func (s *jwtService) ValidateToken(endcodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(endcodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return token, err
	}

	// check status token is valid or not
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return token, errors.New("Token has expired")
		}
	} else {
		return token, errors.New("Invalid token")
	}

	return token, nil
}
