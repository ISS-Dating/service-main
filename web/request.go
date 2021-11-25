package web

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResgisterInfo struct {
}
