package service

import (
	"encoding/json"
	"github.com/cgentry/gus/ecode"
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	_ "github.com/cgentry/gus/storage/sqlite"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

const t_service_test_db = `/tmp/test_service.sqlite`

func generateCaller() *record.User {
	u := record.NewTestUser()
	u.SetDomain(`Test`)
	return u

}
func destroyStore() {
	if t_service_test_db != `:memory:` {
		os.Remove(t_service_test_db)
	}
}

func TestServiceRegister(t *testing.T) {
	sr := NewServiceRegister()

	store, err := storage.Open("sqlite", t_service_test_db, "")
	if err != nil {
		t.Errorf("Error opening store: %s", err.Error())
	}

	defer destroyStore()
	defer store.Close()

	store.CreateStore()
	sr.Client = generateCaller()

	sr.RequestHead = request.NewHead()
	sr.RequestHead.Id = sr.Client.LoginName
	sr.RequestHead.Domain = sr.Client.Domain

	reg := request.NewRegister()

	sr.RequestBody = reg

	Convey("Send Bad Requests in", t, func() {
		sr.UserStore = store
		sr.Client = generateCaller()

		pack := sr.Run(sr)
		rtnHead := pack.Head.(response.Head)
		So(rtnHead.Message, ShouldContainSubstring, ecode.ErrMissingName.Error())
		So(rtnHead.Code, ShouldEqual, ecode.ErrMissingName.Code())
		So(pack.IsBodySet(), ShouldBeFalse)
		So(pack.Body, ShouldBeBlank)

		reg.Name = "name"
		sr.RequestBody = reg
		pack = sr.Run(sr)
		rtnHead = pack.Head.(response.Head)
		So(rtnHead.Message, ShouldContainSubstring, ecode.ErrMissingPassword.Error())

	})

	Convey("Send Simple register request in", t, func() {

		reg.Login = "*Login"
		reg.Name = "*TestSimpleRegister"
		reg.Email = "johndoe@golang.go"
		reg.Password = "12345678abcdefg"

		sr.RequestBody = reg
		sr.Client = generateCaller()

		pack := sr.Run(sr)
		rtnHead := pack.Head.(response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

		userRtn := record.UserReturn{}
		err := json.Unmarshal([]byte(pack.Body), &userRtn)
		So(err, ShouldBeNil)
		So(userRtn.LoginName, ShouldEqual, reg.Login)
		So(userRtn.FullName, ShouldEqual, reg.Name)
		So(userRtn.Email, ShouldEqual, reg.Email)

		// DUPLICATE EMAIL ERROR
		pack = sr.Run(sr)

		rtnHead = pack.Head.(response.Head)
		So(rtnHead.Message, ShouldEqual, ecode.ErrDuplicateEmail.Error())
		So(rtnHead.Code, ShouldEqual, ecode.ErrDuplicateEmail.Code())

		So(len(pack.Body), ShouldEqual, 0) // No data when an error occurs

	})
	Convey("Simple login/logout", t, func() {

		sl := NewServiceLogin()
		sl.UserStore = store

		reqLogin := request.NewLogin()
		reqLogin.Login = "*Login"
		reqLogin.Password = "12345678abcdefg"

		sl.RequestBody = reqLogin
		sl.Client = generateCaller()

		pack := sl.Run(sl)
		rtnHead := pack.Head.(response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)

		userRtn := record.UserReturn{}
		err := json.Unmarshal([]byte(pack.Body), &userRtn)
		So(err, ShouldBeNil)
		So(userRtn.LoginName, ShouldEqual, reqLogin.Login)
		So(userRtn.FullName, ShouldEqual, `*TestSimpleRegister`)
		So(userRtn.Email, ShouldEqual, `johndoe@golang.go`)

		so := NewServiceLogout()
		so.UserStore = store
		so.Client = generateCaller()

		reqLogout := request.NewLogout()
		reqLogout.Token = userRtn.Token
		so.RequestBody = reqLogout

		pack = so.Run(so)

		rtnHead = pack.Head.(response.Head)
		So(rtnHead.Message, ShouldBeBlank)
		So(rtnHead.Code, ShouldEqual, 200)
		So(pack.IsBodySet(), ShouldBeFalse)
		So(pack.Body, ShouldBeBlank)

		pack = so.Run(so)
		rtnHead = pack.Head.(response.Head)
		So(rtnHead.Message, ShouldEqual, ecode.ErrUserNotLoggedIn.Error())
		So(rtnHead.Code, ShouldEqual, ecode.ErrUserNotLoggedIn.Code())
		So(pack.IsBodySet(), ShouldBeFalse)
		So(pack.Body, ShouldBeBlank)
	}) /*

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

			pack := ServiceLogin(ctrl, caller, p)
			rtnHead := pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldNotBeBlank)
			So(rtnHead.Message, ShouldEqual, ecode.ErrUserNotFound.Error())
			So(rtnHead.Code, ShouldEqual, ecode.ErrUserNotFound.Code())

			So(pack.IsBodySet(), ShouldBeFalse)
			So(pack.Body, ShouldBeBlank)

			p.SetBody(`{Login: 10:21:55}`)
			pack = ServiceLogin(ctrl, caller, p)
			rtnHead = pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldNotBeBlank)
			So(rtnHead.Message, ShouldEqual, ecode.ErrInvalidBody.Error())
			So(rtnHead.Code, ShouldEqual, ecode.ErrInvalidBody.Code())

			So(pack.IsBodySet(), ShouldBeFalse)
			So(pack.Body, ShouldBeBlank)

		})

		Convey("Bad logout", t, func() {

			h := request.NewHead()
			h.Domain = `Test`
			h.Id = `ID`
			p := record.NewPackage()

			p.SetHead(h)
			p.SetSecret([]byte(`secret`))
			p.SetBody(`{Token: 10:21:55}`)

			So(p.GetSignature(), ShouldNotEqual, "")

			pack := ServiceLogout(ctrl, caller, p)

			rtnHead := pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldNotBeBlank)
			So(rtnHead.Message, ShouldEqual, ecode.ErrInvalidBody.Error())
			So(rtnHead.Code, ShouldEqual, ecode.ErrInvalidBody.Code())

			So(pack.IsBodySet(), ShouldBeFalse)
			So(pack.Body, ShouldBeBlank)

		})
		Convey(`Test Updating user`, t, func() {

			reqLogin := request.NewLogin()
			reqLogin.Login = "*Login"
			reqLogin.Password = "12345678abcdefg"

			originalUserRecord, err := ctrl.DataStore.UserFetch(caller.Domain, storage.FIELD_LOGIN, `*Login`)
			So(err, ShouldBeNil)

			h := request.NewHead()
			h.Domain = `Test`
			h.Id = `ID`

			p := record.NewPackage()
			p.SetSecret([]byte(`secret`))
			p.SetHead(h)
			p.SetBody(reqLogin)

			So(p.GetSignature(), ShouldNotEqual, "")

			pack := ServiceLogin(ctrl, caller, p)
			rtnHead := pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldBeBlank)
			So(rtnHead.Code, ShouldEqual, 200)
			userRtn := record.UserReturn{}
			err = json.Unmarshal([]byte(pack.Body), &userRtn)
			So(err, ShouldBeNil)
			token := userRtn.Token

			// NOW UPDATE USING TOKEN
			reqUpdate := request.NewUpdate()
			reqUpdate.Login = "*LoginNew"
			reqUpdate.Name = "*Name New"
			reqUpdate.Email = "*Email New"
			reqUpdate.Token = token

			p.SetBody(reqUpdate)
			options := NewOptions()
			options.Set(PERMIT_LOGIN, true)
			pack = ServiceUpdate(ctrl, caller, p, options)

			rtnHead = pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldEqual, "Fields updated: Login")
			So(rtnHead.Code, ShouldEqual, 200)
			userRtn2 := record.UserReturn{}
			err = json.Unmarshal([]byte(pack.Body), &userRtn2)
			So(err, ShouldBeNil)
			So(userRtn2.LoginName, ShouldEqual, reqUpdate.Login)
			So(userRtn2.FullName, ShouldEqual, userRtn.FullName)
			So(userRtn2.Email, ShouldEqual, userRtn.Email)

			//
			// Attempt number 3: email only
			options = NewOptions()
			options.Set(PERMIT_EMAIL, true)
			reqUpdate.Login = "*Login Old"
			p.SetBody(reqUpdate)
			pack = ServiceUpdate(ctrl, caller, p, options)

			rtnHead = pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldEqual, "Fields updated: Email")
			So(rtnHead.Code, ShouldEqual, 200)
			userRtn3 := record.UserReturn{}
			err = json.Unmarshal([]byte(pack.Body), &userRtn3)
			So(err, ShouldBeNil)

			So(userRtn3.LoginName, ShouldEqual, userRtn2.LoginName)
			So(userRtn3.FullName, ShouldEqual, userRtn2.FullName)
			So(userRtn3.Email, ShouldEqual, reqUpdate.Email)

			//
			// Attempt number 4: ALL of them (except password)
			options = NewOptions()
			options.Set(PERMIT_ALL, true)
			reqUpdate.Login = "*Login all"
			reqUpdate.Email = "email@all.com"
			reqUpdate.Name = "*Name all"

			p.SetBody(reqUpdate)
			pack = ServiceUpdate(ctrl, caller, p, options)

			rtnHead = pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldEqual, "Fields updated: Login, Name, Email")
			So(rtnHead.Code, ShouldEqual, 200)
			userRtn4 := record.UserReturn{}
			err = json.Unmarshal([]byte(pack.Body), &userRtn4)
			So(err, ShouldBeNil)

			So(userRtn4.LoginName, ShouldEqual, reqUpdate.Login)
			So(userRtn4.FullName, ShouldEqual, reqUpdate.Name)
			So(userRtn4.Email, ShouldEqual, reqUpdate.Email)

			//
			// Attempt number 5: Password
			options = NewOptions()
			options.Set(PERMIT_PASSWORD, true)
			reqUpdate.Login = "*Login all"
			reqUpdate.Email = "email@all.com"
			reqUpdate.Name = "*Name all"
			reqUpdate.OldPassword = `12345678abcdefg`
			reqUpdate.NewPassword = `abcdefg987654321`

			p.SetBody(reqUpdate)
			pack = ServiceUpdate(ctrl, caller, p, options)

			rtnHead = pack.Head.(*response.Head)
			So(rtnHead.Message, ShouldEqual, "Fields updated: Password")
			So(rtnHead.Code, ShouldEqual, 200)
			userRtn5 := record.UserReturn{}
			err = json.Unmarshal([]byte(pack.Body), &userRtn5)
			So(err, ShouldBeNil)

			So(userRtn5.LoginName, ShouldEqual, userRtn4.LoginName)
			So(userRtn5.FullName, ShouldEqual, userRtn4.FullName)
			So(userRtn5.Email, ShouldEqual, userRtn4.Email)

			lastUserRecord, err := ctrl.DataStore.UserFetch(caller.Domain, storage.FIELD_LOGIN, reqUpdate.Login)
			So(err, ShouldBeNil)
			So(lastUserRecord.Password, ShouldNotEqual, originalUserRecord.Password)
		})*/
}
