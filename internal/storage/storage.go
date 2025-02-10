package storage

import (
	"errors"
)

var (
	ErrContainerNotFound = errors.New("st: container not found")
	//ErrContainerNotFound = errors.New("container not found")
)

const (
	dbPostgres = "postgres"
	dbMysql    = "mysql"
)

//func  (cfg *config.Config) {
//	switch  {
//
//	}
//}
