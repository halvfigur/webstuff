package webstuff

import "errors"

var ErrBadRequest = errors.New("bad request")
var ErrUnauthorized = errors.New("unauthorized")
var ErrExpired = errors.New("expired")
