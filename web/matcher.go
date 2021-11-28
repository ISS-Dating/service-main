package web

import (
	"encoding/json"
	"net/http"

	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/service"
)

type Matcher struct {
	Canteen *Canteen
}

func NewMatcher(s service.Interface) *Matcher {
	return &Matcher{
		Canteen: NewCanteen(s),
	}
}

func (m *Matcher) match(w http.ResponseWriter, req *http.Request) {
	enableCors(&w, req)
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}
	callback := m.Canteen.AddUser(&user)
	friend := <-callback

	json.NewEncoder(w).Encode(friend)
}

func (m *Matcher) answer(w http.ResponseWriter, req *http.Request) {
	enableCors(&w, req)
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var answers model.Questionary
	json.NewDecoder(req.Body).Decode(&answers)

	callback := m.Canteen.GetAnswer(&user, &answers)
	replyAnswers := <-callback

	json.NewEncoder(w).Encode(replyAnswers)
}

func (m *Matcher) status(w http.ResponseWriter, req *http.Request) {
	enableCors(&w, req)
	user, status := auth(req)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	var stats Status
	json.NewDecoder(req.Body).Decode(&stats)

	callback := m.Canteen.GetStatus(&user, &stats)
	replyStats := <-callback

	json.NewEncoder(w).Encode(replyStats)
}

func (m *Matcher) Start() {
	go m.Canteen.Poll()

	http.HandleFunc("/match", m.match)
	http.HandleFunc("/answer", m.answer)
	http.HandleFunc("/status", m.status)

	http.ListenAndServe(":8091", nil)
}
