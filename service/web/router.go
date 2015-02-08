// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"fmt"
	"github.com/cgentry/gus/record/configure"
	"net/http"
)

type ServiceHandler struct {
	config *configure.Configure
}

func NewService(c *configure.Configure) *ServiceHandler {

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) { httpRegister(c, w, r) })
	//http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.request) { ServiceLogin(c, w, r) })
	//http.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.request) { ServiceLogout(c, w, r) })
	//http.HandleFunc("/authenticate/", func(w http.ResponseWriter, r *http.request) { ServiceUpdate(c, w, r) })
	//http.HandleFunc("/enable/", func(w http.ResponseWriter, r *http.request) { ServiceEnable(c, w, r) })
	//http.HandleFunc("/disable/", func(w http.ResponseWriter, r *http.request) { ServiceDisable(c, w, r) })

	return &ServiceHandler{config: c}
}

func (s *ServiceHandler) Start() {
	serviceAddress := fmt.Sprintf("%s:%d", c.Service.Host, c.Service.Port)
	http.ListenAndServe(serviceAddress, nil)
}

/*
func ServiceLogout(c *gofig.Configuration, w http.ResponseWriter, r *http.request) {
	return
}
func ServiceUpdate(c *gofig.Configuration, w http.ResponseWriter, r *http.request) {
	return
}
func ServiceEnable(c *gofig.Configuration, w http.ResponseWriter, r *http.request) {
	return
}
func ServiceDisable(c *gofig.Configuration, w http.ResponseWriter, r *http.request) {
	return
}

func sendError(c *gofig.Configuration, w http.ResponseWriter) {

}
*/
