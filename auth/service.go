package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// AUTH USING JWT
// HOW TO GENERATE TOKEN
// VALIDATE TOKEN
type Service interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

var SECRET = []byte("BWASTARTUP_s3cr3t")

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userId int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(encodedToken *jwt.Token) (interface{}, error) {
		_, ok := encodedToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}
