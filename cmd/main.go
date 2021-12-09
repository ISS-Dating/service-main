package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/ISS-Dating/service-main/migrate"
	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/repo"
	"github.com/ISS-Dating/service-main/service"
	"github.com/ISS-Dating/service-main/web"
)

var (
	localSetup = false
)

func main() {
	var db *sql.DB
	var err error
	if localSetup {
		db, err = sql.Open("postgres", "host=localhost user=postgres password=12345 port=5432 dbname=q_date sslmode=disable")
	} else {
		var stop bool
		for !stop {
			time.Sleep(time.Second * 8)
			db, _ = sql.Open("postgres", "host=postgres user=postgres password=12345 port=5432 dbname=postgres sslmode=disable")
			err = db.Ping()
			if err != nil {
				log.Println("Error connecting to db")
			} else {
				stop = true
			}
		}
	}
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect to db: ", err.Error())
	}

	err = migrate.ApplyMigrations(db, "/server/connections/database")
	if err != nil {
		log.Fatal("Can't apply migrations: ", err.Error())
	}

	service := service.NewService(repo.NewRepo(db))
	server := web.NewServer(service)
	matcher := web.NewMatcher(service)

	service.Repo.CreateUser(model.User{Username: adminLogin, Password: adminPassword, Role: model.RoleAdministrator})

	go matcher.Start()
	server.Start()
}

var adminLogin = "admin"
var adminPassword = "admin"
