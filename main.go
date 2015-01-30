package main

import (
	"github.com/cgentry/gus/command"
	"path/filepath"
	"io/ioutil"
	// Here is where you include all the Store/Drivers you want
	_ "github.com/cgentry/gus/storage/sqlite"
	// _ "github.com/cgentry/gus/storage/mysql"
	// _ "github.com/cgentry/gus/storage/mongo"

	// Here is where you include all the Encryption/Drivers you want
	_ "github.com/cgentry/gus/encryption/drivers/plaintext"
	_ "github.com/cgentry/gus/encryption/drivers/bcrypt"
	_ "github.com/cgentry/gus/encryption/drivers/sha512"
)

const (
	DEFAULT_CONFIG_FILENAME = "/etc/gus/config.json"
	DEFAULT_CONFIG_PERMISSIONS = 0600
)

var configFileName string
func main(){

}

func GetConfigFileName()( Filename string , DirExists error, FileExists error ){
	if configFileName == "" {
		configFileName = DEFAULT_CONFIG_FILENAME
	}
	_,DirExists = os.Stat( filepath.Dir( configFilename ) )
	_,FileExists = os.Stat( configFileName )
	Filename = configFileName
	return
}

func GetConfigFile() ( []byte , error ){
	fname, _,_ := GetConfigFileName()
	return ioutil.ReadFile( fname )
}

func SaveConfigFile( jsonString []byte ) error {
	file, direrror, _ := GetConfigFileName()
	if direrror{
		return direrror
	}
	return ioutil.WriteFile( file, jsonString, DEFAULT_CONFIG_PERMISSIONS )
}
