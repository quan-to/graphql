package gqlerrors

import (
	"errors"

	"github.com/quan-to/graphql/language/location"
)

type ExtendedError interface {
	error
	Extensions() map[string]interface{}
}

type FormattedError struct {
	Message    string                    `json:"message"`
	Locations  []location.SourceLocation `json:"locations"`
	Path       []interface{}             `json:"path,omitempty"`
	Extensions map[string]interface{}    `json:"extensions,omitempty"`
}

func (g FormattedError) Error() string {
	return g.Message
}

func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error) FormattedError {
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		ret := FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Path:      err.Path,
		}
		if err := err.OriginalError; err != nil {
			if extended, ok := err.(ExtendedError); ok {
				ret.Extensions = extended.Extensions()
			}
		}
		return ret
	case Error:
		return FormatError(&err)
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
		}
	}
}

func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := make([]FormattedError, 0)
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}

func FormatErrorsFunc(f func (err error) FormattedError, errs ...error) []FormattedError {
	formattedErrors := make([]FormattedError, 0)
	for _, err := range errs {
		formattedErrors = append(formattedErrors, f(err))
	}
	return formattedErrors
}