package Return

import (
	"time"
)

//
// Interface
// ---------------------------------------------------------------------------------------------------- //
type Interface interface {
	ReturnClear()
	ReturnSetPrefix(format any, args ...any)
	ReturnSetError(format any, args ...any)
	ReturnAddError(format string, args ...any)
	ReturnGetError() error
	ReturnSetWarning(format any, args ...any)
	ReturnAddWarning(format any, args ...any)
	ReturnGetWarning() error
	ReturnGetTime() time.Time
	ReturnIsError() bool
	ReturnIsNotError() bool
	ReturnIsWarning() bool
	ReturnIsNotWarning() bool
	ReturnPrint()

	// Error
}

func NewInterface() Interface {
	return new(Error)
}

func (e *Error) ReturnClear() {
	e.Clear()
}

func (e *Error) ReturnSetPrefix(format any, args ...any) {
	e.SetPrefix(format, args...)
}

func (e *Error) ReturnSetError(format any, args ...any) {
	e.SetError(format, args...)
}

func (e *Error) ReturnAddError(format string, args ...any) {
	e.AddError(format, args...)
}

func (e *Error) ReturnGetError() error {
	return e.GetError()
}

func (e *Error) ReturnSetWarning(format any, args ...any) {
	e.SetWarning(format, args...)
}

func (e *Error) ReturnAddWarning(format any, args ...any) {
	e.AddWarning(format, args...)
}

func (e *Error) ReturnGetWarning() error {
	return e.GetWarning()
}

func (e *Error) ReturnGetTime() time.Time {
	return e.GetTime()
}

func (e *Error) ReturnIsError() bool {
	return e.IsError()
}

func (e *Error) ReturnIsNotError() bool {
	return e.IsNotError()
}

func (e *Error) ReturnIsWarning() bool {
	return e.IsWarning()
}

func (e *Error) ReturnIsNotWarning() bool {
	return e.IsNotWarning()
}

func (e *Error) ReturnPrint() {
	e.Print()
}
