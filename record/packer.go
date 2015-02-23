package record

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"github.com/cgentry/gus/record/head"
	"reflect"
)

type BodyInterface interface {
	Check() error
}

type Packer interface {
	SetBodyMarshal(interface{}) error
	SetBody(string)
	SetHead(head.HeaderInterface)
	SetSecret([]byte)
	ClearSecret()

	GetHead() head.HeaderInterface
	GetBody() string
	GetSecret() []byte

	IsPackageComplete() bool
	IsHeadSet() bool
	IsBodySet() bool
}

// Compute signature creates an HMAC signed signature, using sha256, of the body of
// the request. The request is a JSON encoded string and the secret must be held in
// the key. The key should be the secret key of the ID stored in the header.
func computeSignature(p Packer) (sig []byte) {
	secret := p.GetSecret()
	if secret != nil {
		mac := hmac.New(sha256.New, secret)
		mac.Write([]byte(p.GetBody()))
		sig = mac.Sum(nil)
	}
	return sig
}

// Check to see if the signature in the header is valid.
func GoodSignature(p Packer) bool {
	secret := p.GetSecret()
	if nil != secret && p.GetHead().IsSignatureSet() {
		sig, err := p.GetHead().GetSignature()
		if err == nil {
			return hmac.Equal(sig, computeSignature(p))
		}
	}

	return false
}

// SignPackage with a base64-encoded HMAC of the body contents.
func SignPackage(p Packer) {
	p.GetHead().SetSignature(computeSignature(p))
	return
}

// Package describes what is coming from a request. It is similar to, but different from,
// the response package. It conforms to the Packer interface
type Package struct {
	Version  float32
	BodyType string

	Head   head.HeaderInterface
	Body   string
	code   int
	secret []byte
}

// NewPackage creates a new package and returns the address to the caller.
func NewPackage() Packer {
	return &Package{
		Head:     head.New(),
		Body:     "{}",
		Version:  1.0,
		BodyType: "empty",
	}
}

// Check to see if the package is complete
func (p *Package) IsPackageComplete() bool {
	return p.IsBodySet() && p.IsHeadSet()
}

// CheckHead will determine if we have all the data required for a Head in the package
func (p *Package) IsHeadSet() bool {
	return p.Head != nil && p.Head.IsTimeSet()
}

// IsBodySet checks to see if the body is blank or not
func (p *Package) IsBodySet() bool {
	return p.Body != ""
}

// Return the HEAD of the package
func (p *Package) GetHead() head.HeaderInterface {
	return p.Head
}

// SetHead will copy the head to our head, then sign the package if it is complete.
func (p *Package) SetHead(head head.HeaderInterface) {
	p.Head = head
	return
}

// GetBody returns the body string
func (p *Package) GetBody() string {
	return p.Body
}

func (p *Package) SetBodyMarshal(body interface{}) (err error) {
	var byteBody []byte
	p.BodyType = reflect.TypeOf(body).String()
	byteBody, err = json.Marshal(body)
	if err == nil {
		p.Body = string(byteBody)
	}
	return
}

// SetBody returns the value that is in the body. If the package is complete, it will be signed.
func (p *Package) SetBody(body string) {
	p.Body = body
	return
}

// GetSecret returns the secret used to sign the package
func (p *Package) GetSecret() []byte {
	return p.secret
}

// SetSecret sets the secret used to sign the package
func (p *Package) SetSecret(s []byte) {
	p.secret = s
	return
}

func (p *Package) ClearSecret() {
	p.secret = nil
	return
}
