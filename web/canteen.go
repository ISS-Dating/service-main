package web

import (
	"sync"

	"github.com/ISS-Dating/service-main/model"
	"github.com/ISS-Dating/service-main/service"
)

const QueueSize = 10000

type Entry struct {
	User             *model.User
	Callback         chan *model.User
	AnswerCallback   chan *model.Questionary
	AnswerCallbackIn chan *model.Questionary
	StatusCallback   chan *model.Status
	StatusCallbackIn chan *model.Status
}

type Canteen struct {
	Queue   chan Entry
	List    sync.Map
	Service service.Interface
}

type Table struct {
	EntryA  *Entry
	EntryB  *Entry
	Service service.Interface
}

func NewCanteen(s service.Interface) *Canteen {
	return &Canteen{
		Queue:   make(chan Entry, QueueSize),
		Service: s,
	}
}

func (c *Canteen) Poll() {
	for {
		var ok bool
		var userA Entry
		var userB Entry
		for !ok {
			userA, ok = <-c.Queue
		}
		ok = false
		for !ok {
			userB, ok = <-c.Queue
		}

		entryA := &Entry{
			User:             userA.User,
			Callback:         make(chan *model.User),
			AnswerCallback:   make(chan *model.Questionary),
			AnswerCallbackIn: make(chan *model.Questionary),
			StatusCallback:   make(chan *model.Status),
			StatusCallbackIn: make(chan *model.Status),
		}
		entryB := &Entry{
			User:             userB.User,
			Callback:         make(chan *model.User),
			AnswerCallback:   make(chan *model.Questionary),
			AnswerCallbackIn: make(chan *model.Questionary),
			StatusCallback:   make(chan *model.Status),
			StatusCallbackIn: make(chan *model.Status),
		}

		userA.Callback <- userB.User
		userB.Callback <- userA.User

		c.List.Store(userA.User.ID, entryA)
		c.List.Store(userB.User.ID, entryB)

		table := NewTable(c.Service, entryA, entryB)
		go table.Poll()
	}
}

func NewTable(s service.Interface, entryA, entryB *Entry) *Table {
	return &Table{
		Service: s,
		EntryA:  entryA,
		EntryB:  entryB,
	}
}

func (c *Canteen) AddUser(user *model.User) chan *model.User {
	callback := make(chan *model.User)
	entry := Entry{
		User:     user,
		Callback: callback,
	}
	c.Queue <- entry
	return callback
}

func (c *Canteen) GetAnswer(user *model.User, answers *model.Questionary) chan *model.Questionary {
	entryRaw, ok := c.List.Load(user.ID)
	if !ok {
		return nil
	}

	entry := entryRaw.(*Entry)
	entry.AnswerCallbackIn <- answers
	return entry.AnswerCallback
}

func (c *Canteen) GetStatus(user *model.User, status *model.Status) chan *model.Status {
	entryRaw, ok := c.List.Load(user.ID)
	if !ok {
		return nil
	}

	entry := entryRaw.(*Entry)
	entry.StatusCallbackIn <- status
	return entry.StatusCallback
}

func (t *Table) Poll() {
	answersA := <-t.EntryA.AnswerCallbackIn
	answersB := <-t.EntryB.AnswerCallbackIn

	t.EntryA.AnswerCallback <- answersB
	t.EntryB.AnswerCallback <- answersA

	statusA := <-t.EntryA.StatusCallbackIn
	statusB := <-t.EntryB.StatusCallbackIn

	t.EntryA.StatusCallback <- statusB
	t.EntryB.StatusCallback <- statusA
}
