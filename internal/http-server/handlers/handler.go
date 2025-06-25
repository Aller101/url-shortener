package handlers

import "errors"

var (
	ErrVoidURL  = errors.New("url is blank")
	ErrLengtURL = errors.New("url lengts is wrong")
)
