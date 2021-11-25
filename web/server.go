package web

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

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

	token, err := createToken(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	embedToken(w, token)
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

	token, err := createToken(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	embedToken(w, token)

	json.NewEncoder(w).Encode(user)
}

// /get_photo endpoint
func (s *Server) getPhoto(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	file, err := os.ReadFile(path.Join("static", user.Username+".png"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.Write(file)
}

// set_photo endpoint
func (s *Server) setPhoto(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	photo, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if len(photo) >= 1024*1024*10 {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Photo is too big, 10Mb max"))
		return
	}

	os.Remove(path.Join("static", user.Username+".png"))
	os.WriteFile(path.Join("static", user.Username+".png"), photo, os.ModePerm)
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
	err := os.MkdirAll("static", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	validateKey = &signKey.PublicKey

	http.HandleFunc("/login", s.login)
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/update", s.update)
	http.HandleFunc("/stats", s.stats)
	http.HandleFunc("/block", s.block)
	http.HandleFunc("/chat_list", s.chatList)

	http.HandleFunc("/set_photo", s.setPhoto)
	http.HandleFunc("/get_photo", s.getPhoto)

	http.ListenAndServe(":8090", nil)
}

func NewServer(service service.Interface) *Server {
	return &Server{
		Service: service,
	}
}
