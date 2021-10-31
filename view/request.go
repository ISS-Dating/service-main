package view

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResgisterInfo struct {
}
