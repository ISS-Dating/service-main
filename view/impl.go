package view

import "net/http"

type DefaultView struct {
	service service.Interface
}

func New() Interface {
	return &DefaultView{}
}

func (d *DefaultView) Login(w http.ResponseWriter, req *http.Request) {
	_, err := d.service.Login()
}
