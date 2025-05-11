package main

import (
	"log"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/231031/authservice"
)

func main() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println(err)
	}
	time.Local = loc

	db := authservice.GetEnv("DB_URL", "")
	str := authservice.GetEnv("PORT", "50001")
	port, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	repository := authservice.NewRepository(db)
	service := authservice.NewService(repository)
	err = authservice.ListenGRPC(service, int(port))
	if err != nil {
		log.Fatal(err)
	}
}
