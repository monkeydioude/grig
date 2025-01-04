package fs

import "errors"

var ErrReadDir = errors.New("error reading directory")
var ErrReadEntryFile = errors.New("unable to read entry file")
var ErrReadIniFile = errors.New("unable to read ini file")
