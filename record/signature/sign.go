package signature

import ()

func New() *Signature {
	return &Signature{ sum : []byte(``)}
}

type Signature struct {
	sum []byte
}

func (s *Signature) SetSignature(newSum []byte) { s.sum = newSum }
func (s *Signature) GetSignature() []byte       { return s.sum }
