package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/ISS-Dating/service-main/repo"
	"github.com/ISS-Dating/service-main/service"
	"github.com/ISS-Dating/service-main/web"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=12345 dbname=q_date sslmode=disable")
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}

	go web.StartStaticServer()
	server := web.NewServer(service.NewService(repo.NewRepo(db)))
	server.Start()
}
