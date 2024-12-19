package errors

import "errors"

// os
var ErrUnknownOS = errors.New("unknown OS")

// fs
var ErrReadDir = errors.New("error reading directory")
var ErrReadEntryFile = errors.New("unable to read entry file")
var ErrReadIniFile = errors.New("unable to read ini file")

// misc
var ErrUnmarshaling = errors.New("unable to read entry file")
