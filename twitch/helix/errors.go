package helix

import "errors"

var (
	Err400                = errors.New("error 400")
	Err404                = errors.New("error 404")
	ErrNotAuthorized      = errors.New("not authorized")
	ErrUnknown            = errors.New("unknown error")
	ErrInvalidResponse    = errors.New("invalid response")
	ErrUnexpectedResponse = errors.New("unexpected response")
	ErrInvalidBody        = errors.New("invalid body")
)
