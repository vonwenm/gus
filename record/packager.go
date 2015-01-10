package record

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type HeaderInterface interface {
	Check() error
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
	secret []byte
}

func NewPackage() *Package {
	return &Package{}
}

func (p *Package) IsPackageComplete() bool {
	return p.IsBodySet() && p.IsHeadSet()
}

func (p *Package) IsBodySet() bool {
	return p.Body != ""
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
		return p, err
	}
	p.Body = string(bodyJson)
	p.setSignature()
	return p, err
}

func (p *Package) SetBodyString(body string) *Package {
	p.Body = body
	return p
}

func (p *Package) setSignature() {
	p.Head.SetSignature(base64.StdEncoding.EncodeToString(p.computeSignature()))
}

func (p *Package) computeSignature() []byte {
	if nil != p.secret && p.IsPackageComplete() {
		mac := hmac.New(sha256.New, p.secret)
		mac.Write([]byte(p.Body))
		return mac.Sum(nil)
	}
	return []byte(``)
}

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

func (p *Package) GetSignature() string {
	return p.Head.GetSignature()
}

func (p *Package) SetSecret(NewSecret []byte) {
	p.secret = NewSecret
}

func (p *Package) ClearSecret() {
	p.secret = nil
}
