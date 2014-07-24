package service

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/cgentry/gus/record"
	"bytes"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func TestAddGet(t * testing.T) {
	sr := NewServiceRequest()
	sr.Add("b", "2").Add("a", "1").Add("c", "3")

	Convey("Values should be there", t, func() {
		k, found := sr.Get("a")
		So(k, ShouldEqual, "1")
		So(found, ShouldBeTrue)
	})
}

func TestSortKeys(t * testing.T) {

	sr := NewServiceRequest()
	sr.Add("b", "2").Add("a", "1").Add("c", "3")

	Convey("Values should be sorted", t, func() {
		order := sr.SortedKeys()
		lastKey := ""
		for _, key := range order {
			So(bytes.Compare([]byte(lastKey), []byte(key)), ShouldBeLessThan, 1)
			lastKey = key
		}
	})
}

func TestParseParms(t * testing.T) {

	req, _ := http.NewRequest("GET", "http://example.com/one/two/three/hhmmaacc?four=4&five=5&six=6", nil)
	slash := []string{ "cmd", "domain", "caller", "hmac"}
	qparm := []string{ "four", "five"}

	Convey("Command variables should be there", t, func() {
		sr, err := ParseParms(req, slash, qparm)
		So(err, ShouldBeNil)

		key, found := sr.Get("cmd")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "one")

		key, found = sr.Get("domain")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "two")

		key, found = sr.Get("caller")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "three")

		key, found = sr.Get("hmac")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "hhmmaacc")

		key, found = sr.Get("four")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "4")

		key, found = sr.Get("five")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "5")

		_, found = sr.Get("six")
		So(found, ShouldBeFalse)

	})
}

func TestParseParms_multiple(t * testing.T) {

	req, _ := http.NewRequest("GET", "http://example.com/one?two=2&two=3&two=4", nil)
	slash := []string{ "cmd" }
	qparm := []string{ "two"}

	Convey("Command variables should be there", t, func() {
		sr, _ := ParseParms(req, slash, qparm)

		key, found := sr.Get("two")
		So(found, ShouldBeTrue)
		So(key, ShouldEqual, "2")

	})
}

func TestParseParms_missing_path(t * testing.T) {

	req, _ := http.NewRequest("GET", "http://example.com/one?two=2&two=3&two=4", nil)
	slash := []string{ "cmd" , "nope"}
	qparm := []string{ "two"}

	Convey("Command variables should be there", t, func() {
		sr, err := ParseParms(req, slash, qparm)

		So(err, ShouldNotBeNil)
		_, found := sr.Get("nope")
		So(found, ShouldBeFalse)

	})
}
/*
 * Make sure we can get a hash and verify the record using an in-line date
 */
func TestCreateHash(t * testing.T) {

	cmd := "register"

	for i := 0; i < 10; i++ {

		secret := record.CreateSalt(50)
		tm	:= time.Now()
		myDate := tm.Format(time.RFC1123)

		h := hmac.New(sha256.New, []byte(secret))        // Start the hmac up
		h.Write([]byte("/" + cmd + "/domain/name"))                        // Adding in the fresh command hash

		h.Write([]byte("pwdsomethingusercharles" + myDate))
		myHmac := base64.StdEncoding.EncodeToString(h.Sum(nil))

		url := fmt.Sprintf(`http://example.com/%s/domain/name?user=charles gentry&pwd=something&date=%s&hmac=%s`,
			cmd, myDate , myHmac )
		req, _ := http.NewRequest("GET", url, nil)

		sr := NewServiceRequest()
		sr.Add("cmd", cmd).Add("user", "charles").Add("pwd", "something").Add("hmac", myHmac)
		sr.SetPathKeys([]string{"cmd" , "hmac"})
		sr.SetQueryKeys( []string{ "pwd" , "user" , "hmac"})

		Convey("Test basic hash creation", t, func() {
			key, err := CreateRestfulHmac(secret, req, &sr)

			So(err, ShouldBeNil)
			So(CompareHmac(key, &sr), ShouldBeTrue)
		})
	}
}

func TestCreateHash_HeaderDate(t * testing.T) {

	cmd := "register"

	for i := 0; i < 10; i++ {

		secret := record.CreateSalt(50)
		tm	:= time.Now()
		myDate := tm.Format(time.RFC1123)

		h := hmac.New(sha256.New, []byte(secret))        // Start the hmac up
		h.Write([]byte("/" + cmd + "/domain/name"))                        // Adding in the fresh command hash

		h.Write([]byte("pwdsomethingusercharles" + myDate))
		myHmac := base64.StdEncoding.EncodeToString(h.Sum(nil))

		url := fmt.Sprintf(`http://example.com/%s/domain/name?user=charles gentry&pwd=something&hmac=%s`,
			cmd , myHmac )
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add( HEADER_DATE , myDate)

		sr := NewServiceRequest()
		sr.Add("cmd", cmd).Add("user", "charles").Add("pwd", "something").Add("hmac", myHmac)
		sr.SetPathKeys([]string{"cmd" , "hmac"})

		sr.SetQueryKeys( []string{ "pwd" , "user" , "hmac"})

		Convey("Test basic hash creation", t, func() {
			key, err := CreateRestfulHmac(secret, req, &sr)

			So(err, ShouldBeNil)
			So(CompareHmac(key, &sr), ShouldBeTrue)
		})
	}
}

/*
 * Make sure we can get a hash and verify the record using an in-line date
 */
func TestCreateHash_embedded_characters(t * testing.T) {

	cmd := "register"

		secret := record.CreateSalt(50)
		tm	:= time.Now()
		myDate := tm.Format(time.RFC1123)


		h := hmac.New(sha256.New, []byte(secret))        // Start the hmac up
		h.Write([]byte("/" + cmd + "/domain/name"))                        // Adding in the fresh command hash

		h.Write([]byte("pwdsomethingusercharles" + myDate))
		myHmac := base64.StdEncoding.EncodeToString(h.Sum(nil))

		url := fmt.Sprintf(`http://example.com/%s/domain/name?user=charles gentry&pwd=something&date=%s&hmac=%s`,
			cmd, myDate , myHmac )
		req, _ := http.NewRequest("GET", url, nil)

		sr := NewServiceRequest()
		sr.Add("cmd", cmd).Add("user", "charles").Add("pwd", "something").Add("hmac", myHmac)
		sr.SetPathKeys([]string{"cmd" , })

	sr.SetQueryKeys( []string{ "pwd" , "user" , "hmac"})

		Convey("Test basic hash creation", t, func() {
			key, err := CreateRestfulHmac(secret, req, &sr)
			So(err, ShouldBeNil)
			So(CompareHmac(key, &sr), ShouldBeTrue)
		})

}


