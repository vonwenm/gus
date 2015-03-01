package signature

import (
	"encoding/base64"
)

func New() *Signature {
	return &Signature{Signature: ``, isSet: false}
}

// Signature contains the base 64-encoded string from the byte array
// The methods attached will always return the []byte arrays originall
// passed in. This aids in serialisation of the signature and consolidates
// the encoding in one place
type Signature struct {
	Signature string
	isSet     bool
}

func (s *Signature) SetSignature(newSum []byte) {
	s.Signature = base64.StdEncoding.EncodeToString(newSum)
	s.isSet = true
}
func (s *Signature) GetSignature() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Signature)
}
func (s *Signature) IsSignatureSet() bool {
	return s.isSet
}
