package authservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	ErrAuthFailed = errors.New("login Failed")
)

type Service interface {
	getAccessToken(authCode string) (*AuthGoogleResp, error)
	getInfo()
	refreshToken()
}

type service struct {
	repo   Repository
	errLog *log.Logger
}

func NewService(r Repository) Service {
	l := log.New(nil, "error from service : ", log.Ldate|log.LUTC|log.Lshortfile)
	return &service{repo: r, errLog: l}
}

func (s *service) getAccessToken(authCode string) (*AuthGoogleResp, error) {
	payload := GoogleInfo{
		Code:        authCode,
		GrantType:   "authorization",
		RedirectURI: os.Getenv("GOOGLE_REDIRECT_URI"),
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		s.errLog.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer(data))
	if err != nil {
		s.errLog.Println(err)
		return nil, err
	}

	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", genBasicAuthHeader(client_id, client_secret))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.errLog.Println(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrAuthFailed
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &AuthGoogleResp{}, nil
}

func (s *service) getInfo()      {}
func (s *service) refreshToken() {}
