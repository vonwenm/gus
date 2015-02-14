// Package web is the interface between HTTP/HTTPS calls and the service layer. It acts as an intermediate
// router to call the correct service routine, depending upon the route called.
//
// How Routing works:
// A route table is created that contains the HTTP Server routine to be called and the Service routine
// to be called. This acts as an interface layer to the services coded. It means that differnt interfaces
// (such as RPC) could be used by writing a different server later. This service layer allows mocking of
// routines by passing in a different map of what routines to call and what http requests can be answered.
package web
