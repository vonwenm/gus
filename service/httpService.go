package service

// This is an interface between the Service layer and the HTTP service. It allows
// Service to be more general and to be called by other interfaces (e.g. RPC)
import (
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/storage"
	"net/http"
)
type ServiceControl struct {
	DataStore *storage.Store
	ClientStore *storage.Store
}

func NewServiceControl() *ServiceControl {
	return &ServiceControl{}
}


func httpRegister(c *configure.Configure, w http.ResponseWriter, r *http.Request) {
	var err error

	ctrl := NewServiceControl()

	if err = ctrl.HttpOpenStore(c, w); err != nil {
		return
	}
	defer ctrl.DataStore.Reset()   // If it holds state...clear it
	defer ctrl.DataStore.Release() // Make sure to release any locks
	defer ctrl.ClientStore.Reset()   // If it holds state...clear it
	defer ctrl.ClientStore.Release() // Make sure to release any locks

	requestPackage, requestHead, err := httpGetBody(r)
	if err != nil {
		httpErrorWrite(w, http.StatusBadRequest, err.Error())
		return
	}
	caller, err := ctrl.HttpFindSystemUser(requestHead.Id)
	if err != nil {
		httpErrorWrite(w, http.StatusUnauthorized, "")
		return
	}

	responsePackage := ServiceRegister(ctrl, caller, requestPackage)
	httpResponseWrite(w, responsePackage)

}
