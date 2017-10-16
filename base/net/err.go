package net

import "errors"

var (
	ErrSizeExceedLimit = errors.New("size exceed limit")
	ErrInvalidType     = errors.New("invalid type")
)
