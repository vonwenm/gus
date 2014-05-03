package main

import ( 
	"github.com/cgentry/gus/service"
	"github.com/cgentry/gofig"
 )

func main(){
	config,err := gofig.NewConfigurationFromIniString("[data]\na=b");
	if err != nil {
		panic( err )
	}
	service.NewService( config );
}
