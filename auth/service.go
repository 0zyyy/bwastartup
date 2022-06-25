package auth

import "github.com/golang-jwt/jwt/v4"

// AUTH USING JWT
// HOW TO GENERATE TOKEN
// VALIDATE TOKEN
type Service interface {
	GenerateToken(userId int) (string, error)
}

type jwtService struct {
}

var SECRET = []byte("BWASTARTUP_s3cr3t")

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userId int) (string, error) {
	claim := jwt.MapClaims{}
	claim["userid"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
