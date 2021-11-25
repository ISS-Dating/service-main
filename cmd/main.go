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

func main() {
	time.Sleep(time.Second * 10)
	db, err := sql.Open("postgres", "host=postgres user=postgres password=12345 port=5432 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}

	server := web.NewServer(service.NewService(repo.NewRepo(db)))
	server.Start()
}
