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
)

func generateCaller() *record.User {
	u := record.NewTestUser()
	u.SetDomain(`Test`)
	return u

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


		fmt.Println( "Call serviceregister")
		pack := ServiceRegister(store, caller, p)
		fmt.Println( "Return serviceregister")
		fmt.Println(pack)

		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

		// DUPLICATE EMAIL ERROR
		pack = ServiceRegister(store, caller, p)
		fmt.Println(pack)

		rtnHead = pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldEqual, storage.ErrDuplicateEmail.Error())
		So(rtnHead.Code, ShouldEqual, storage.ErrDuplicateEmail.Code())

	})
	Convey("Simple login" , t , func() {
		reg := request.NewRegister()
		reg.Login = "*Login"
		reg.Password = "12345678abcdefg"

		h := request.NewHead()
		h.Domain = `Test`
		h.Id = `ID`

		p := record.NewPackage()
		p.SetSecret([]byte(`secret`))
		p.SetHead(h)
		p.SetBody(reg)

		So(p.GetSignature(), ShouldNotEqual, "")

		pack := ServiceLogin(store,caller,p)
		fmt.Println(pack)
		rtnHead := pack.Head.(*response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

	})
}

