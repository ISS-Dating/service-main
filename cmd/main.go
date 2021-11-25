package main

import "github.com/ISS-Dating/service-main/web"

func main() {
	server := web.NewServer()
	server.Start()
}
