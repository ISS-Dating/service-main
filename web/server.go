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

	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/service"
)

type Server struct {
	Service service.Interface
}

// /login endpoint
func (s *Server) login(w http.ResponseWriter, req *http.Request) {
	var login genericRequest
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
	var login genericRequest
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
	_, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var requestData genericRequest
	err := json.NewDecoder(req.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	file, err := os.ReadFile(path.Join("static", requestData.Username+".png"))
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
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var reqUser model.User
	err := json.NewDecoder(req.Body).Decode(&reqUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if reqUser.Username != user.Username {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	reqUser.Role = user.Role

	user, err = s.Service.UpdateUser(reqUser)
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

// /stats endpoint
func (s *Server) statsByUsername(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var data genericRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Role == model.RoleAdministrator || user.Role == model.RoleModerator {
		searchUser, err := s.Service.GetUserByUsername(data.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.NewEncoder(w).Encode(searchUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		return
	}

	w.WriteHeader(http.StatusForbidden)
}

// /ban endpoint
func (s *Server) ban(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var data genericRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Service.BanUser(user, data.Username, data.Mod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

// /mod endpoint
func (s *Server) mod(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var data genericRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Service.ModUser(user, data.Username, data.Mod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

// /friends endpoint
func (s *Server) friends(w http.ResponseWriter, req *http.Request) {
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	friends, err := s.Service.ListFriends(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	matched := MatchedList{
		Friends: friends,
	}

	json.NewEncoder(w).Encode(matched)
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
	http.HandleFunc("/stats_username", s.statsByUsername)
	http.HandleFunc("/ban", s.ban)
	http.HandleFunc("/mod", s.mod)

	http.HandleFunc("/set_photo", s.setPhoto)
	http.HandleFunc("/get_photo", s.getPhoto)

	http.HandleFunc("/friends", s.friends)

	http.ListenAndServe(":8090", nil)
}

func NewServer(service service.Interface) *Server {
	return &Server{
		Service: service,
	}
}
