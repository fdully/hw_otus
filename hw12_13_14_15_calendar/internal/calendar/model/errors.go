package model

import "errors"

var (
	ErrDateBusy     = errors.New("this time is busy")
	ErrAlreadyExist = errors.New("this object already exists")
	ErrNotExist     = errors.New("this object not exist")
)
