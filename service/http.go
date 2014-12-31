package service

import (
	"net/http"
	//"github.com/cgentry/gus/storage"
	"encoding/json"
	"errors"
	"github.com/cgentry/gus/record"
	"github.com/cgentry/gus/record/request"
	"github.com/cgentry/gus/record/response"
	"github.com/cgentry/gus/storage"
	"io/ioutil"
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
func httpFindSystemUser(caller_guid string) (*record.User, error) {

	// We need the calling system's secret. This is the token for the caller
	drive := storage.GetDriver()
	caller, err := drive.FetchUserByGuid(caller_guid)
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
