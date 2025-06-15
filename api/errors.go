package api

import "errors"

var (
	ErrWrongCredentials = errors.New("WRONG_CREDENTIALS")
	ErrFetchError       = errors.New("FETCH_ERROR")
	ErrParseError       = errors.New("PARSE_ERROR")

	ErrUnknownError = errors.New("UNKNOWN_ERROR")
)
