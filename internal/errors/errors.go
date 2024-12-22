package errors

import "errors"

// service
// service/os
var ErrUnknownOS = errors.New("unknown OS")

// service/fs
var ErrReadDir = errors.New("error reading directory")
var ErrReadEntryFile = errors.New("unable to read entry file")
var ErrReadIniFile = errors.New("unable to read ini file")

// http
var ErrHttpUnknownInternalServerError = errors.New("unknown internal server error")

// misc
var ErrUnmarshaling = errors.New("unable to unmarshal data")
var ErrMarshaling = errors.New("unable to marshal data")
var ErrReadingFile = errors.New("unable to read file")
