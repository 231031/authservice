package authservice

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	getUserbySub(sub string) (*Auth, error)
	getUserbyId(id string) (*Auth, error)
	createUser(user *Auth) (*Auth, error)
	updateRefresh(id, refresh_token string, expire_in time.Time) error
}

type repository struct {
	db     *gorm.DB
	errLog *log.Logger
}

func NewRepository(url string) Repository {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(nil, "error from db :", log.Ldate|log.LUTC|log.Lshortfile)

	return &repository{db: db, errLog: l}
}

func (r *repository) createUser(authInfo *Auth) (*Auth, error) {
	resp := r.db.Create(authInfo)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil, resp.Error
	}
	return authInfo, nil
}

func (r *repository) getUserbySub(sub string) (*Auth, error) {
	authInfo := &Auth{}
	resp := r.db.Where(&Auth{Sub: sub}).First(authInfo)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil, resp.Error
	}
	return authInfo, nil
}

func (r *repository) getUserbyId(id string) (*Auth, error) {
	authInfo := &Auth{}
	resp := r.db.Where(&Auth{Id: id}).First(authInfo)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil, resp.Error
	}
	return authInfo, nil
}

func (r *repository) updateRefresh(id, refresh_token string, expire_in time.Time) error {
	resp := r.db.Model(&Auth{Id: id}).Updates(
		Auth{RefreshToken: refresh_token, ExpiresIn: expire_in},
	)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil
	}
	return nil
}
