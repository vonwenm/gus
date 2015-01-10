// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//
package service

import (
	"github.com/cgentry/gofig"
	"net/http"
)

type ServiceHandler struct{}

func NewService(c *gofig.Configuration) {

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) { httpRegister(c, w, r) })
	//http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) { ServiceLogin(c, w, r) })
	//http.HandleFunc("/logout/", func(w http.ResponseWriter, r *http.Request) { ServiceLogout(c, w, r) })
	//http.HandleFunc("/authenticate/", func(w http.ResponseWriter, r *http.Request) { ServiceUpdate(c, w, r) })
	//http.HandleFunc("/enable/", func(w http.ResponseWriter, r *http.Request) { ServiceEnable(c, w, r) })
	//http.HandleFunc("/diable/", func(w http.ResponseWriter, r *http.Request) { ServiceDisable(c, w, r) })

	//http.ListenAndServe(":8181", nil)
}

/*
func ServiceLogout(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceUpdate(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceEnable(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}
func ServiceDisable(c *gofig.Configuration, w http.ResponseWriter, r *http.Request) {
	return
}

func sendError(c *gofig.Configuration, w http.ResponseWriter) {

}
*/
