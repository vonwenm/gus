package record

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type HeaderInterface interface {
	IsTimeSet() bool
	GetSignature() string
	SetSignature(string)
}

type BodyInterface interface {
	Check() error
}

type PackagerInterface interface {
	SetBody(interface{}) (*Package, error)
	SetBodyString(string) *PackagerInterface
	SetHead(*HeaderInterface) *PackagerInterface
	SetSecret([]byte) *PackagerInterface
	GoodSignature() bool
	IsPackageComplete() bool
	IsHeadSet() bool
	IsBodySet() bool
}

type Package struct {
	Head   HeaderInterface
	Body   string
	code   int
	secret []byte
}

func NewPackage() *Package {
	return &Package{}
}

func (p *Package) IsPackageComplete() bool {
	return p.IsBodySet() && p.IsHeadSet()
}

func (p *Package) IsBodySet() bool {
	return p.Body != ``
}

func (p *Package) IsHeadSet() bool {
	return p.Head != nil && p.Head.IsTimeSet()
}

func (p *Package) SetHead(h HeaderInterface) *Package {
	p.Head = h
	p.setSignature()
	return p
}

func (p *Package) SetBody(body interface{}) (*Package, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		p.Body = ``
		return p, err
	}
	p.Body = string(bodyJson)
	p.setSignature()
	return p, err
}

func (p *Package) SetBodyString(body string) *Package {
	p.Body = body
	p.setSignature()
	return p
}

// Compute a new signature for the body and save it in the header. The signature is a base64
// encoded string.
func (p *Package) setSignature() {
	if p.IsHeadSet() {
		sig := p.computeSignature()
		p.Head.SetSignature(base64.StdEncoding.EncodeToString(sig))
	}
}

// Compute signature creates an HMAC signed signature, using sha256, of the body of
// the request. The request is a JSON encoded string and the secret must be held in
// the key. The key should be the secret key of the ID stored in the header.
func (p *Package) computeSignature() []byte {
	if p.IsPackageComplete() && nil != p.secret {
		mac := hmac.New(sha256.New, p.secret)
		mac.Write([]byte(p.Body))
		return mac.Sum(nil)
	}
	return []byte(``)
}

// Check to see if the signature in the header is valid.
func (p *Package) GoodSignature() bool {
	if nil != p.secret {
		if sig := p.Head.GetSignature(); sig != "" {
			if headSignature, err := base64.StdEncoding.DecodeString(sig); err == nil {
				return hmac.Equal(headSignature, p.computeSignature())
			}
		}
	}

	return false
}

// Return the signature of the header. This is a simple convenience wrapper
func (p *Package) GetSignature() string {
	return p.Head.GetSignature()
}

// Set the package secret for encoding/decoding purposes
func (p *Package) SetSecret(NewSecret []byte) {
	p.secret = NewSecret
}

// Clear the secret from the package. As this is a private (lowercase) variable
// the value will never be encoded.
func (p *Package) ClearSecret() {
	p.secret = nil
}
