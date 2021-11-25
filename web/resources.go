package web

import (
	"log"
	"net/http"
)

func StartStaticServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	log.Println("Static files server started")
	http.ListenAndServe(":8091", nil)
}
