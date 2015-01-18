package service

import (
	"encoding/json"
	"errors"
	"github.com/cgentry/gofig"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	"io/ioutil"
	"net/http"
)

func httpGetBody(r *http.Request) (requestPackage *record.Package, requestHead request.Head, err error) {
	requestPackage = record.NewPackage()
	requestHead = request.NewHead()
	var OK bool

	httpRequestBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(httpRequestBody, &requestPackage)
		if err == nil {
			requestHead, OK = requestPackage.Head.(request.Head)
			if !OK {
				err = errors.New("Invalid request head")
			}
		}
	}

	return
}

// Given a GUID, find the user's record in the database. Only system users
// will be returned. If none are found, return an error
func httpFindSystemUser(store *storage.Store, callerGuid string) (*record.User, error) {
	caller, err := store.FetchUserByGuid(callerGuid)
	if err != nil || !caller.IsSystem {
		return nil, errors.New("Invalid user id or password")
	}

	return caller, nil
}

func httpErrorWrite(w http.ResponseWriter, code int, msg string) {
	responseHead := response.NewHead()
	responseHead.Code = code
	responseHead.Message = msg
	responsePackage := record.NewPackage()
	responsePackage.SetHead(responseHead)
	httpResponseWrite(w, responsePackage)
}

func httpResponseWrite(w http.ResponseWriter, responsePackage *record.Package) {
	responseHead, OK := responsePackage.Head.(response.Head)
	if !OK {
		responseHead = response.NewHead()
	}
	httpResponseBody, _ := json.Marshal(responsePackage)
	w.Write(httpResponseBody)
	http.Error(w, responseHead.Message, responseHead.Code)
}

func httpOpenStore(c *gofig.Configuration, w http.ResponseWriter) (ds *storage.Store, err error) {
	var driverName, driverDsn string

	if driverName, err = c.GetString(`store`, `driver`); err == nil {
		if driverDsn, err = c.GetString(`store`, `dsn`); err == nil {
			options := c.GetStringWithDefault(`store`, `options`, ``)
			ds, err = storage.Open(driverName, driverDsn, options)
		}
	}
	if err != nil {
		httpErrorWrite(w, 500, err.Error())
		return
	}
	return
}
