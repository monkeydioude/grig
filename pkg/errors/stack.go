package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Custom error type with stack trace
type StackError struct {
	msg   string
	cause error
	stack []uintptr
}

// NewStackError creates a new StackError with a stack trace
func NewStackError(msg string, cause error) *StackError {
	stack := make([]uintptr, 32)
	length := runtime.Callers(3, stack[:]) // Skip 3 frames to ignore runtime.Callers and helpers
	return &StackError{
		msg:   msg,
		cause: cause,
		stack: stack[:length],
	}
}

// Error returns the error message
func (e *StackError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.cause)
	}
	return e.msg
}

// Unwrap allows compatibility with errors.Is and errors.As
func (e *StackError) Unwrap() error {
	return e.cause
}

// StackTrace formats the captured stack trace
func (e *StackError) StackTrace() string {
	var sb strings.Builder
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

// Implement fmt.Formatter for custom %+v formatting.
// This should deprecated tho, as we never know when a new flag might appear.
func (e *StackError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		// Handle %+v for stack trace
		if f.Flag('+') {
			fmt.Fprintf(f, "%s\n", e.Error())
			fmt.Fprintf(f, "Stack trace:\n%s", e.StackTrace())
			return
		}
		fmt.Fprintf(f, "%s", e.Error())
	case 's':
		fmt.Fprintf(f, "%s", e.Error())
	case 'q':
		fmt.Fprintf(f, "%q", e.Error())
	}
}

// Wrap adds a stack trace and wraps an error with a message
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return NewStackError(msg, err)
}

// Wrapf adds a stack trace and wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Errorf(format, args...).Error()
	return NewStackError(msg, err)
}

// Print display the error with the stack trace.
// Should be used insteaf of fmt.Printf("%+v", err)
func Print(err error) {
	fmt.Println(err)
	if se, ok := err.(*StackError); ok {
		fmt.Println("Stack trace:")
		fmt.Println(se.StackTrace())
	}
}
