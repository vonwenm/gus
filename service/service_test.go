package service

import (
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/storage/mock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"fmt"
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
	_ "github.com/cgentry/gus/storage/sqlite"
	"encoding/json"

)

func generateCaller() *record.User {
	u := record.NewTestUser()
	u.SetDomain(`Test`)
	return u

}

func TestBadRegister( t * testing.T){
	caller := generateCaller()
	store, _ := storage.Open("sqlite", ":memory:")
	defer store.Close()
	store.CreateStore()

	Convey("Send Bad Requests in", t, func() {

		pack := ServiceRegister(store,caller,nil)
		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldContainSubstring, storage.ErrInvalidHeader.Error())
		So(rtnHead.Code, ShouldEqual, storage.ErrInvalidHeader.Code())
		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)

		p := record.NewPackage()
		p.Head = request.NewHead()
		pack = ServiceRegister(store,caller,p)
		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldContainSubstring, `No domain`)
		So(rtnHead.Code, ShouldEqual, storage.ErrInvalidHeader.Code())
		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)

		h := request.NewHead()
		h.Domain = `Test`
		h.Id = `ID`
		p.Head = h

		pack = ServiceRegister(store,caller,p)
		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldContainSubstring, `Invalid Checksum`)
		So(rtnHead.Code, ShouldEqual, storage.ErrInvalidHeader.Code())
		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)
	})
}
func TestSimpleRegister(t *testing.T) {

	mock.RegisterMockStore()
	caller := generateCaller()
	store, err := storage.Open("sqlite", ":memory:")
	if err != nil {
		t.Error("Culd not open database: " + err.Error())
	}
	defer store.Close()
	store.CreateStore()

	Convey("Send Simple register request in", t, func() {
		reg := request.NewRegister()
		reg.Login = "*Login"
		reg.Name = "*FullName"
		reg.Email = "johndoe@golang.go"
		reg.Password = "12345678abcdefg"

		h := request.NewHead()
		h.Domain = `Test`
		h.Id = `ID`

		p := record.NewPackage()
		p.SetSecret([]byte(`secret`))
		p.SetHead(h)
		p.SetBody(reg)

		So(p.GetSignature(), ShouldNotEqual, "")

		fmt.Println("Call serviceregister")
		pack := ServiceRegister(store, caller, p)
		fmt.Println("\nReturn serviceregister\n")
		fmt.Println(pack)

		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

		userRtn := record.UserReturn{}
		err := json.Unmarshal( []byte(pack.Body), &userRtn )
		So( err, ShouldBeNil)
		So( userRtn.LoginName, ShouldEqual, reg.Login )
		So( userRtn.FullName, ShouldEqual, reg.Name)
		So( userRtn.Email, ShouldEqual, reg.Email)

		// DUPLICATE EMAIL ERROR
		pack = ServiceRegister(store, caller, p)
		fmt.Println("\nReturn bad:\n")
		fmt.Println(pack)

		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldEqual, storage.ErrDuplicateEmail.Error())
		So(rtnHead.Code, ShouldEqual, storage.ErrDuplicateEmail.Code())

		So( len( pack.Body), ShouldEqual, 0 )	// No data when an error occurs


	})
	Convey("Simple login/logout", t, func() {
		reqLogin := request.NewLogin()
		reqLogin.Login = "*Login"
		reqLogin.Password = "12345678abcdefg"

		h := request.NewHead()
		h.Domain = `Test`
		h.Id = `ID`

		p := record.NewPackage()
		p.SetSecret([]byte(`secret`))
		p.SetHead(h)
		p.SetBody(reqLogin)

		So(p.GetSignature(), ShouldNotEqual, "")

		pack := ServiceLogin(store, caller, p)
		fmt.Println(pack)
		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

		userRtn := record.UserReturn{}
		err := json.Unmarshal( []byte(pack.Body), &userRtn )
		So( err, ShouldBeNil)
		So( userRtn.LoginName, ShouldEqual, reqLogin.Login )
		So( userRtn.FullName, ShouldEqual, `*FullName`)
		So( userRtn.Email, ShouldEqual, `johndoe@golang.go`)

		reqLogout := request.NewLogout()
		reqLogout.Token = userRtn.Token
		p.SetBody(reqLogout)
		pack = ServiceLogout(store,caller,p)

		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)
		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)

		pack = ServiceLogout(store,caller,p)
		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldEqual, storage.ErrUserNotLoggedIn.Error())
		So(rtnHead.Code, ShouldEqual, storage.ErrUserNotLoggedIn.Code())
		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)
	})

	Convey("Bad login", t, func() {
		reg := request.NewRegister()
		reg.Login = "*LoginXX"
		reg.Password = "12345678abcdefg"

		h := request.NewHead()
		h.Domain = `Test`
		h.Id = `ID`

		p := record.NewPackage()
		p.SetSecret([]byte(`secret`))
		p.SetHead(h)
		p.SetBody(reg)

		So(p.GetSignature(), ShouldNotEqual, "")

		pack := ServiceLogin(store, caller, p)
		fmt.Println(pack)
		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldNotBeBlank)
		So(rtnHead.Message, ShouldEqual, storage.ErrUserNotFound.Error())
		So(rtnHead.Code, ShouldEqual, storage.ErrUserNotFound.Code())

		So(pack.IsBodySet(), ShouldBeFalse )
		So(pack.Body, ShouldBeBlank)

	})
}
