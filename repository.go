package authservice

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	getUserbySub(sub string) (*User, error)
	createUser(user *User) (*User, error)
	updateRefreshToken(id, refresh_token string) error
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

func (r *repository) createUser(user *User) (*User, error) {
	id := uuid.New().String()
	user.id = id

	resp := r.db.Create(user)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil, resp.Error
	}
	return user, nil
}

func (r *repository) getUserbySub(sub string) (*User, error) {
	user := &User{}
	resp := r.db.Where(&User{sub: sub}).First(user)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil, resp.Error
	}
	return user, nil
}

func (r *repository) updateRefreshToken(id, refresh_token string) error {
	resp := r.db.Model(&User{id: id}).Update("refresh_token", refresh_token)
	if resp.Error != nil {
		r.errLog.Println(resp.Error)
		return nil
	}
	return nil
}
