package errors

import "errors"

// service/ini
var ErrReadIniFile = errors.New("unable to read ini file")

// html
var ErrEmptyLinkText = errors.New("WithNav: element.Link.Text can absolutely not be empty. WithNav can derive a page's name from element.Link.Href, unless it's empty or a single '/'")
var ErrEmptyLinkHref = errors.New("WithNav: element.Link.Href cannot be empty")

// model
var ErrModelVerifyInvalidValue = errors.New("invalid value during verification")

// misc
var ErrUnmarshaling = errors.New("unable to unmarshal data")
var ErrMarshaling = errors.New("unable to marshal data")
var ErrReadingFile = errors.New("unable to read file")
var ErrCheckingFile = errors.New("unable to check file")
var ErrWritingFile = errors.New("unable to write to file")
var ErrCreatingFile = errors.New("unable to create file")
var ErrNilPointer = errors.New("nil pointer")
