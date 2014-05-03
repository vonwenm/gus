package service

import (
	"net/http"
	"github.com/cgentry/gofig"
	"fmt"
	"encoding/json"
	"errors"
	"strings"
	"sort"
)


type ServiceRequest map[string]string
type ServiceHandler struct {}

type StatusReturn struct {
	Status	int
	Message	string
}


var stdPathParam = []string{"cmd","domain","caller","hmac"}
func NewService( c * gofig.Configuration ){

	http.HandleFunc( "/register/" , func(w http.ResponseWriter, r *http.Request){ServiceRegister( c , w, r )} )
	/*

	*/

	http.ListenAndServe(":8181" , nil  )
}

func ( sr * ServiceRequest ) SortedKeys() []string {
	keys := make( []string , len( *sr ))
	i := 0
	for key,_ := range *sr {
		keys[i] = key
		i++
	}
	sort.Strings( keys )
	return keys
}
// ServiceRegister will handle the calling of Registration for the user.
func ServiceRegister( c *gofig.Configuration , w http.ResponseWriter, r *http.Request )  {
	// Need the request params. Since we have a standard format, parse by default
	qparam := []string{ "email", "login","name", "password"}
	srequest,err := ParseParms( r , stdPathParam, qparam )
	if err != nil {
		ReturnError( w , CODE_BAD_CALL, err )
	}else{
		fmt.Fprintf(w , "Good boy!<br>%s<br>" , srequest )
	}

}

// ParseParms takes the url path and puts it into a standard map
// the request always looks like:
// /cmd/domain/caller-appid/hmac/identifier....
//			login: identifier(login-name)/password
//			auth:  identifier (token)
//			logout: identifier
//			register: identifier(login-name)/password/email
//			lookup: identifier/type (email|name|guid)
//			save:   identifier/type(session|user)/name(of item)
//			retrieve: identifier/type/name
//
//			inactive: indentifer/type
//			active:
//
func ParseParms(r *http.Request , list []string , qparam []string ) (ServiceRequest , error) {
	sr := make( ServiceRequest )

	parts := strings.Split( r.URL.Path , "/")[1:]
	plen := len( parts )
	if plen != len( list ) {
		return sr , errors.New("Path was invalid")
	}

	// Match up all the keys
	for i,key := range list {
		sr[key] = parts[i]
	}

	// AND now for the parameters (follows the ? portion)
	query := r.URL.Query()
	for _,key := range qparam {
		if _,found := query[key]; ! found {
			return sr, errors.New("Missing query parameter '" + key +"'")
		}else{
		sr[ key ] = query.Get( key )
		}
	}
	return sr , nil
}


func ReturnError( w http.ResponseWriter , code int , err error) {

	msg := StatusReturn{ Status : code , Message : err.Error() }
	rtn,_ := json.Marshal( msg )
	w.WriteHeader(http.StatusBadRequest)
	w.Write( rtn )
}

