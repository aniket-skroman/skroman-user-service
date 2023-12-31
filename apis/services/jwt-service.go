package services

import (
	"fmt"
	"os"

	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(userID, userType, dept string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID         string    `json:"user_id"`
	UserType       string    `json:"user_type"`
	UserDepartment string    `json:"dept"`
	CreatedAt      time.Time `json:"created_at"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "ydhnwb",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretkey := os.Getenv("JWT_SECRET")

	if secretkey != "" {
		secretkey = "ydhnwb"
	}
	return secretkey
}

func (j *jwtService) GenerateToken(UserID, userType, dept string) string {

	claims := &jwtCustomClaim{
		UserID,
		userType,
		dept,
		time.Now().Add(10 * time.Minute),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}

	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
