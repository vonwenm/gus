package service

import (
	"github.com/cgentry/gus/storage"
)

type ServiceControl struct {
	DataStore *storage.Store
}

func NewServiceControl() *ServiceControl {
	return &ServiceControl{}
}
