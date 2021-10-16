package web

import (
	"net/http"
)

// /login endpoint
func login(w http.ResponseWriter, req *http.Request) {
}

// /register endpoint
func register(w http.ResponseWriter, req *http.Request) {
}

// /update endpoint
func update(w http.ResponseWriter, req *http.Request) {
}

// /stats endpoint
func stats(w http.ResponseWriter, req *http.Request) {
}

// /block endpoint
func block(w http.ResponseWriter, req *http.Request) {
}

// /chat_list endpoint
func chatList(w http.ResponseWriter, req *http.Request) {
}

func StartHttpServer() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/update", update)
	http.HandleFunc("/stats", stats)
	http.HandleFunc("/block", block)
	http.HandleFunc("/chat_list", chatList)

	http.ListenAndServe(":8090", nil)
}
