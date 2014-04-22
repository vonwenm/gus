package memory

import (
	"testing"
	"fmt"
	"github.com/cgentry/gus/storage"
	"github.com/cgentry/gus"
)


func TestStart( t *testing.T ){
	fmt.Println( " Storage is " + storage.GetDriverName() )
	//startUp()
	fmt.Println("OK")
}

func TestRegister( t * testing.T ){
	drive := storage.GetDriver()

	user := gus.NewTestUser()
	drive.RegisterUser( &user )

	drive.Find
}
