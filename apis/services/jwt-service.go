package services

import (
	"encoding/json"
	"errors"
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
	s := time.Now()
	claims := &jwtCustomClaim{
		UserID,
		userType,
		dept,
		time.Now().Add(10 * time.Minute),
		jwt.StandardClaims{
			ExpiresAt: int64(5 * time.Minute),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	et := time.Since(s)
	fmt.Println("time taken for generate token : ", et.Milliseconds())
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		if t.Valid {
			fmt.Println("token valid ", t)
		}
		byteArr, _ := json.Marshal(t.Claims)

		tokenData := new(jwtCustomClaim)

		json.Unmarshal(byteArr, &tokenData)

		d := time.Since(tokenData.CreatedAt)

		if d.Hours() > 24 {
			return nil, errors.New("token has been expired")
		}

		return []byte(j.secretKey), nil
	})
}
