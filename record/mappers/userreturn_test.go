package mappers

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/cgentry/gus/record/tenant"
	"github.com/cgentry/gus/record/response"
	"testing"
    _ 	"github.com/cgentry/gus/encryption/drivers/plaintext"
	"github.com/cgentry/gus/encryption"
	//"time"
)

func TestUserReturn(t *testing.T) {
	encryption.Select("plaintext")
	user := tenant.NewTestUser()
	user.GenerateGuid()
	user.Token = user.CreateToken()
	user.FullName = `FullName`
	user.LoginName = `LoginName`
	user.SetPassword(`ThisIsThePassword`)
	err := user.Login(`ThisIsThePassword`)
	Convey("Setup and check copy function", t, func() {
		So(err, ShouldBeNil)
		rtn := ResponseFromUser( response.NewUserReturn() , user )

		So(rtn.Guid, ShouldEqual, user.Guid)
		So(rtn.Token, ShouldEqual, user.Token)
		So(rtn.FullName, ShouldEqual, user.FullName)
		So(rtn.LoginName, ShouldEqual, user.LoginName)
		So(rtn.Email, ShouldEqual, user.Email)
		So(rtn.LoginAt.Equal(user.LoginAt), ShouldBeTrue)
		So(rtn.LastAuthAt.Equal(user.LastAuthAt), ShouldBeTrue)
		So(rtn.TimeoutAt.Equal(user.TimeoutAt), ShouldBeTrue)
		So(rtn.MaxSessionAt.Equal(user.MaxSessionAt), ShouldBeTrue)
		So(rtn.CreatedAt.Equal(user.CreatedAt), ShouldBeTrue)

	})
}
