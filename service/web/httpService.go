package service

// This is an interface between the Service layer and the HTTP service. It allows
// Service to be more general and to be called by other interfaces (e.g. RPC)
import (
	"github.com/cgentry/gus/record/configure"
	"github.com/cgentry/gus/storage"
	"net/http"
)

func httpRegister(c *configure.Configure, w http.ResponseWriter, r *http.Request) {
	var err error

	ctrl := NewServiceControl( c )

	if err = ctrl.HttpOpenStore(c, w); err != nil {
		return
	}
	defer ctrl.DataStore.Reset()   // If it holds state...clear it
	defer ctrl.DataStore.Release() // Make sure to release any locks
	defer ctrl.ClientStore.Reset()   // If it holds state...clear it
	defer ctrl.ClientStore.Release() // Make sure to release any locks

	ctrl.Request, s.Head, err := httpGetBody(r)
	if err != nil {
		httpErrorWrite(w, http.StatusBadRequest, err.Error())
		return
	}


	responsePackage := ctrl.Register(requestHead, requestPackage)
	httpResponseWrite(w, responsePackage)

}
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
func (s *ServiceControl )HttpFindSystemUser(callerGuid string) (*record.User, error) {
	caller, err := s.ClientStore.FetchUserByGuid(callerGuid)
	if err != nil || !caller.IsSystem {
		return nil, errors.New("Invalid user id or password")
	}
	s.ClientStore.Release()
	s.ClientStore.Reset()
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

func (s *ServiceControl ) httpOpenStore(c *configure.Configure, w http.ResponseWriter) ( err error) {
	var driverName, driverDsn string

	s.DataStore ,err = storage.Open( c.Store.Name, c.Store.Dsn, c.Store.Options)
	if err != nil {
		httpErrorWrite(w, 500, err.Error())
		return
	}
	if c.Service.ClientStore {
		s.ClientStore,err = storage.Open( c.Client.Name, c.Client.Dsn, c.Client.Options)
	}else{
		s.DataStore,err = storage.Open( c.Store.Name, c.Store.Dsn, c.Store.Options)
	}
	if err != nil {
		httpErrorWrite(w, 500, err.Error())
	}
	return
}
