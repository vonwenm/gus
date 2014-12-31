package service

import (
	"github.com/cgentry/gofig"
	"net/http"
)

func httpRegister(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {

	requestPackage, requestHead, err := httpGetBody(r)
	if err != nil {
		httpErrorWrite(w, http.StatusBadRequest, err.Error())
		return
	}
	caller, err := httpFindSystemUser(requestHead.Id)
	if err != nil {
		httpErrorWrite(w, http.StatusUnauthorized, "")
		return
	}
	responsePackage := ServiceRegister(caller, requestPackage)
	httpResponseWrite(w, responsePackage)

}
