package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("RrbrdW1VouBEL2sAd7sa")

type Claims struct {
	UserId   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Level    int    `json:"level"`
	jwt.StandardClaims
}

func GenerateToken(userId int64, username, password string, level int) (string, error) {
	claims := Claims{
		UserId:   userId,
		Username: username,
		Password: password,
		Level:    level,
		StandardClaims: jwt.StandardClaims{
			Issuer: "console",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
