package authservice

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrNotFoundKey = errors.New("env key not found")
)

func convertExpiresTime(expires_in int) time.Time {
	cur_time := time.Now()
	return cur_time.Add(time.Second * time.Duration(expires_in))
}

func genBasicAuthHeader(clientID, clientSecret string) string {
	auth := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func generateToken(id, access_token, role string) (string, error) {
	secretKey := GetEnv("SECRET_KEY", "")
	if secretKey == "" {
		return "", ErrNotFoundKey
	}

	claims := jwt.MapClaims{
		"id":          id,
		"acces_token": access_token,
		"role":        role,
		"exp":         time.Now().Add(time.Hour * 3).Unix(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func getInfoToken(token string) TokenInfo {
	return TokenInfo{}
}

func GetEnv(key, d string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Println("not found key : ", key)
		return d
	}
	return v
}
