package web

import (
	"encoding/json"
	"net/http"

	"github.com/ISS-Dating/service-main/service"
)

type Server struct {
	Service service.Interface
}

// /login endpoint
func (s *Server) login(w http.ResponseWriter, req *http.Request) {
	var login LoginInfo
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := s.Service.Login(login.Username, login.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(user)
}

// /register endpoint
func (s *Server) register(w http.ResponseWriter, req *http.Request) {
	var login LoginInfo
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := s.Service.Register(login.Username, login.Password, login.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(user)
}

// /update endpoint
func (s *Server) update(w http.ResponseWriter, req *http.Request) {
}

// /stats endpoint
func (s *Server) stats(w http.ResponseWriter, req *http.Request) {
}

// /block endpoint
func (s *Server) block(w http.ResponseWriter, req *http.Request) {
}

// /chat_list endpoint
func (s *Server) chatList(w http.ResponseWriter, req *http.Request) {
}

func (s *Server) Start() {
	http.HandleFunc("/login", s.login)
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/update", s.update)
	http.HandleFunc("/stats", s.stats)
	http.HandleFunc("/block", s.block)
	http.HandleFunc("/chat_list", s.chatList)

	http.ListenAndServe(":8090", nil)
}

func NewServer(service service.Interface) *Server {
	return &Server{
		Service: service,
	}
}
