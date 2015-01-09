package request

type Authenticate struct {
	Token    string
}

func NewAuthenticate() *Authenticate {
	return &Authenticate{}
}
