package brreg

import "errors"

var (
	ErrNotFound              = errors.New("this is not the response you are looking for")
	ErrUnmarshallingResponse = errors.New("you best check your struct, fool")
)
