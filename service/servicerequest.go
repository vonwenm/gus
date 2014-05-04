// Copyright 2014 Charles Gentry. All rights reserved.
// Please see the license included with this package
//

package service

import (
	"sort"
)

/*
 *	A ServiceRequest is a simple map that contains all the request parameters as key/value pairs
 *
 */
type ServiceRequest struct {
	Parameters	map[string]string
	PathKeys    []string
}

func NewServiceRequest() ServiceRequest {
	return ServiceRequest{ Parameters : make(map[string]string) }
}

func (sr * ServiceRequest) SetPathKeys( paths []string ) * ServiceRequest {
	sr.PathKeys = paths
	return sr
}

func (sr * ServiceRequest) GetPathKeys( ) []string  {
	return sr.PathKeys
}

// Return all of the keys from the ServiceRequest in a sorted array.
func (sr * ServiceRequest) SortedKeys() []string {
	keys := make([]string, len(sr.Parameters))
	i := 0
	for key, _ := range sr.Parameters {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// Add a new key and value to the service map. Use of this will protect against structure changes
func (sr * ServiceRequest) Add(key, value string) * ServiceRequest {
	sr.Parameters[key] = value
	return sr
}

// Find a value in the service map. Use of this will protect against structure changes
func (sr * ServiceRequest) Get(key string) ( string , bool ) {
	val,found := sr.Parameters[key]
	return val,found
}



