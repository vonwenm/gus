package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/cgentry/gus/storage/mock"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
)

func generateCaller() * record.User {
	u := record.NewTestUser()
	return u

}
func TestSimpleRegister(t *testing.T) {

	mock := mock.RegisterMockStore()
	caller := generateCaller()

	Convey("Send Simple register request in", t, func() {
		reg := request.NewRegister()
		req.Login = "login"
		req.Name = "name"
		req.email = "johndoe@golang.go"
		req.Password = "12345678abcdefg"
		So( req.Check(), ShouldBeTrue )

		pack := ServiceRegister( caller , request )
	})
}
