package errors

import "errors"

// service/ini
var ErrReadIniFile = errors.New("unable to read ini file")

// html
var ErrEmptyLinkText = errors.New("WithNav: element.Link.Text can absolutely not be empty. WithNav can derive a page's name from element.Link.Href, unless it's empty or a single '/'")
var ErrEmptyLinkHref = errors.New("WithNav: element.Link.Href cannot be empty")

// services
var ErrServicesFilepathExists = errors.New("Provided filepath already exists")
var ErrServicesInvalidFilepath = errors.New("Provided filepath does not exist")
var ErrServicesUnableFileParsing = errors.New("Unable to read the provided filepath")
var ErrServicesServicesUpdateFail = errors.New("Could not update services")
var ErrServicesServicesRestartFail = errors.New("Could not restart service")

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
var ErrEmptySlice = errors.New("empty slice")
var ErrConfigSaveFail = errors.New("Could not save config")
var ErrInvalidProvidedParameters = errors.New("Invalid provided parameters")
