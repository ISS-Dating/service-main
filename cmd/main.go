package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/ISS-Dating/service-main/repo"
	"github.com/ISS-Dating/service-main/service"
	"github.com/ISS-Dating/service-main/web"
)

var (
	localSetup = true
)

func main() {
	var db *sql.DB
	var err error
	if localSetup {
		db, err = sql.Open("postgres", "host=localhost user=postgres password=12345 port=5432 dbname=q_date sslmode=disable")
	} else {
		time.Sleep(time.Second * 5)
		db, err = sql.Open("postgres", "host=postgres user=postgres password=12345 port=5432 dbname=postgres sslmode=disable")
	}
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}

	service := service.NewService(repo.NewRepo(db))
	server := web.NewServer(service)
	matcher := web.NewMatcher(service)
	go matcher.Start()
	server.Start()
}
