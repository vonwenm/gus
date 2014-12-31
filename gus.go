package main

import (
	"github.com/cgentry/gofig"
	"github.com/cgentry/gus/service"
)

func main() {
	config, err := gofig.NewConfigurationFromIniString("[data]\na=b")
	if err != nil {
		panic(err)
	}
	service.NewService(config)
}
