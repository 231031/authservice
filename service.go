package authservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAuthFailed   = errors.New("login Failed")
	ErrUnauthorized = errors.New("token expired")
	ErrServerError  = errors.New("internal server error")
)

type Service interface {
	loginAccount(auth_code, role string) (string, error)
	checkAccessToken(token string) (string, error)
	getAccessToken(authCode string) (*AuthGoogleResp, error)
	getInfo(access_token string) (*UserInfoResponse, error)
	refreshToken(id string) error
	updateRefresh(id, refresh_token string, expires_in int) error
}

type service struct {
	repo   Repository
	errLog *log.Logger
}

func NewService(r Repository) Service {
	l := log.New(nil, "error from service : ", log.Ldate|log.LUTC|log.Lshortfile)
	return &service{repo: r, errLog: l}
}

func (s *service) loginAccount(auth_code, role string) (string, error) {
	authResp, err := s.getAccessToken(auth_code)
	if err != nil {
		return "", err
	}

	userResp, err := s.getInfo(authResp.AccessToken)
	if err != nil {
		return "", err
	}

	// check existing user
	user, err := s.repo.getUserbySub(userResp.Sub)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			id := uuid.New().String()
			user.Id = id
			user.Sub = userResp.Sub
			user.Role = role
			user.RefreshToken = authResp.RefreshToken

			_, err = s.repo.createUser(user)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		// update refresh token and expire time
		err = s.updateRefresh(user.Id, authResp.RefreshToken, authResp.ExpiresIn)
		if err != nil {
			return "", ErrServerError
		}
	}

	token, err := generateToken(user.Id, authResp.AccessToken, user.Role)
	if err != nil {
		s.errLog.Fatal(err)
	}

	return token, nil
}

func (s *service) checkAccessToken(token string) (string, error) {
	info := getInfoToken(token)
	_, err := s.getInfo(info.AccessToken)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			err = s.refreshToken(info.Id)
			if err != nil {
				return "", err
			}
		}

	}

	return token, nil
}

func (s *service) getAccessToken(authCode string) (*AuthGoogleResp, error) {
	endpoint := GetEnv("GOOGLE_REDIRECT_URI", "")
	payload := AuthPayload{
		Code:        authCode,
		GrantType:   "authorization",
		RedirectURI: endpoint,
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	endpoint = GetEnv("OAUTH", "https://oauth2.googleapis.com") + "/token"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrAuthFailed
	}

	client_id := GetEnv("CLIENT_ID", "")
	client_secret := GetEnv("CLIENT_SECRET", "")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", genBasicAuthHeader(client_id, client_secret))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrAuthFailed
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	authResp := &AuthGoogleResp{}
	err = json.Unmarshal(data, authResp)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	return authResp, nil
}

func (s *service) getInfo(access_token string) (*UserInfoResponse, error) {
	api := GetEnv("GOOGLE_API_INFO", "")
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.errLog.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		s.errLog.Println(resp)
		return nil, ErrUnauthorized
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	userInfo := &UserInfoResponse{}
	err = json.Unmarshal(respByte, userInfo)
	if err != nil {
		s.errLog.Println(err)
		return nil, ErrServerError
	}

	return userInfo, nil
}

func (s *service) refreshToken(id string) error {
	authInfo, err := s.repo.getUserbyId(id)
	if err != nil {
		return ErrServerError
	}

	payload := RefreshTokenPayload{
		GrantType:    "refresh_token",
		RefreshToken: authInfo.RefreshToken,
	}
	data, err := json.Marshal(&payload)
	if err != nil {
		s.errLog.Println(err)
		return ErrServerError
	}

	endpoint := GetEnv("OAUTH", "https://oauth2.googleapis.com") + "/oauth/refresh"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		s.errLog.Println(err)
		return ErrServerError
	}

	client_id := GetEnv("CLIENT_ID", "")
	client_secret := GetEnv("CLIENT_SECRET", "")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", genBasicAuthHeader(client_id, client_secret))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.errLog.Println(err)
		return ErrAuthFailed
	}
	if resp.StatusCode != http.StatusOK {
		s.errLog.Println(resp)
		return ErrAuthFailed
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.errLog.Println(err)
		return ErrServerError
	}

	authResp := &AuthGoogleResp{}
	err = json.Unmarshal(data, authResp)
	if err != nil {
		s.errLog.Println(err)
		return ErrServerError
	}

	err = s.updateRefresh(authInfo.Id, authInfo.RefreshToken, authResp.ExpiresIn)
	if err != nil {
		return ErrServerError
	}

	return nil
}

func (s *service) updateRefresh(id, refresh_token string, expires_in int) error {
	exp := convertExpiresTime(expires_in)
	err := s.repo.updateRefresh(id, refresh_token, exp)
	if err != nil {
		return err
	}

	return nil
}
