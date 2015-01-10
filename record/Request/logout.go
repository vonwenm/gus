package request

type Logout struct {
	Token string
}

func NewLogout() *Logout {
	return &Logout{}
}
