package request

type Login struct {
	Login    string
	Password string
}

func NewLogin() *Login {
	return &Login{}
}
