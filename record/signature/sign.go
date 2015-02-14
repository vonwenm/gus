package signature

import (

)

func New() * Signature {
	return &Signature{}
}

type Signature struct {
	sum string
}

func (s *Signature) SetSignature(newSum string) { s.sum = newSum }
func (s *Signature) GetSignature() string       { return s.sum }
