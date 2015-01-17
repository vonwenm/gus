package service

import (
	"github.com/cgentry/gofig"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus/storage/mock"
	"net/http"
)

func httpRegister(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {

	mock.RegisterMockStore()
	ctrl := NewServiceControl()
	ctrl.DataStore, _ = storage.Open("mock", ":memory:")

	requestPackage, requestHead, err := httpGetBody(r)
	if err != nil {
		httpErrorWrite(w, http.StatusBadRequest, err.Error())
		return
	}
	caller, err := httpFindSystemUser(ctrl.DataStore, requestHead.Id)
	if err != nil {
		httpErrorWrite(w, http.StatusUnauthorized, "")
		return
	}

	responsePackage := ServiceRegister(ctrl, caller, requestPackage)
	httpResponseWrite(w, responsePackage)

}
