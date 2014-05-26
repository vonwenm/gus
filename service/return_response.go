package service


import (
	"net/http"
	"encoding/json"
	"github.com/cgentry/gus/record"
	"time"
)

func ReturnError(w http.ResponseWriter , code int , err error) {

	msg := record.Status{ Code : code , Message : err.Error() }
	rtn, _ := json.Marshal(msg)
	 ReturnString(w , http.StatusBadRequest , rtn )
}

func ReturnUserJson( w http.ResponseWriter , code int , rtn *record.UserReturn ){
	rtnJson,_ := json.Marshal(rtn)
	 ReturnString(w,code,rtnJson)
}
func ReturnString( w http.ResponseWriter , code int , rtn []byte ){
	now := time.Now().Format(time.RFC3339)
	w.Header().Add( "Content-Type" , "application/json")
	w.Header().Add( HEADER_DATE , now )
	// Create the hmac for the response
	w.WriteHeader(code)
	w.Write( rtn )

}
