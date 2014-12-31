package record

import (
	"crypto/hmac"
	"crypto/sha256"
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
	SetBody(*interface{}) (*Package, error)
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

func (p *Package) SetBody(body *interface{}) (*Package, error) {
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
	p.Head.SetSignature(p.computeSignature())
}

func (p *Package) computeSignature() string {
	if nil != p.secret && p.IsPackageComplete() {
		mac := hmac.New(sha256.New, p.secret)
		mac.Write([]byte(p.Body))
		return string(mac.Sum(nil))
	}
	return ""
}

func (p *Package) GoodSignature() bool {
	headSignature := []byte(p.Head.GetSignature())
	if nil != p.secret && nil == headSignature {
		return hmac.Equal(headSignature, []byte(p.computeSignature()))
	}
	return false
}

func (p *Package) GetSignature() string {
	return p.Head.GetSignature()
}

func (p *Package) SetSecret(NewSecret []byte) {
	p.secret = NewSecret
}
