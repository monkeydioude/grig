package fs

import "errors"

var ErrReadDir = errors.New("error reading directory")
var ErrReadEntryFile = errors.New("unable to read entry file")
var ErrReadIniFile = errors.New("unable to read ini file")
var ErrCheckingFile = errors.New("unable to check file")
var ErrWritingFile = errors.New("unable to write to file")
var ErrCreatingFile = errors.New("unable to create file")
