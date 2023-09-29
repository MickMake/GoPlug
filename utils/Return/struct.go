package Return

import (
	"errors"
	"fmt"
	"os"
	"time"
)

//
// Error
// ---------------------------------------------------------------------------------------------------- //
type Error struct {
	prefix  string
	when    time.Time
	err     error
	warning error
}

var Ok Error

func New() Error {
	return Ok
}

func NewWithPrefix(format string, args ...any) Error {
	var err Error
	err.SetPrefix(format, args...)
	return err
}

func NewError(format any, args ...any) Error {
	var e Error
	e.SetError(format, args...)
	return e
}

func NewWarning(format any, args ...any) Error {
	var e Error
	e.SetWarning(format, args...)
	return e
}

func (e *Error) format(format any, args ...any) string {
	var str string

	switch v := format.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		str = fmt.Sprintf("Integer: %v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		str = fmt.Sprintf("Float64: %v", v)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		str = fmt.Sprintf("%v", v)
		str = fmt.Sprintf(str, args...)
	case error:
		str = fmt.Sprintf("%s", format)
	case Error:
		str = fmt.Sprintf("%s", format)
	}

	return str
}

// Clear - Clears all structure data.
func (e *Error) Clear() {
	*e = Error{
		prefix:  e.prefix,
		when:    time.Time{},
		err:     nil,
		warning: nil,
	}
}

// Copy - Copies an Error structure, (keeping the current prefix).
func (e *Error) Copy(err Error) Error {
	*e = Error{
		prefix:  e.prefix,
		when:    err.when,
		err:     err.err,
		warning: err.warning,
	}
	return *e
}

func (e *Error) SetPrefix(format any, args ...any) {
	str := e.format(format, args...)
	if str == "" {
		return
	}
	e.prefix = str
	e.Clear()
}

func (e *Error) SetError(format any, args ...any) {
	str := e.format(format, args...)
	if str == "" {
		return
	}
	e.when = time.Now()
	e.err = errors.New(str)
	e.warning = nil
}

func (e *Error) AddError(format string, args ...any) {
	str := e.format(format, args...)
	if str == "" {
		return
	}
	e.when = time.Now()
	e.warning = nil
	if e.err == nil {
		e.err = errors.New(str)
		return
	}
	e.err = errors.New(fmt.Sprintf("%s / %s", e.err, str))
}

func (e *Error) GetError() error {
	if e.err == nil {
		return nil
	}
	return errors.New(fmt.Sprintf("%s%v", e.prefix, e.err))
}

func (e *Error) SetWarning(format any, args ...any) {
	str := e.format(format, args...)
	if str == "" {
		return
	}
	e.when = time.Now()
	e.err = nil
	e.warning = errors.New(str)
}

func (e *Error) AddWarning(format any, args ...any) {
	str := e.format(format, args...)
	if str == "" {
		return
	}
	e.when = time.Now()
	e.err = nil
	if e.warning == nil {
		e.warning = errors.New(str)
		return
	}
	e.warning = errors.New(fmt.Sprintf("%v / %s", e.warning, str))
}

func (e *Error) GetWarning() error {
	if e.warning == nil {
		return nil
	}
	return errors.New(fmt.Sprintf("%s%v", e.prefix, e.warning))
}

func (e *Error) GetTime() time.Time {
	return e.when
}

func (e *Error) IsError() bool {
	if e.err == nil {
		return false
	}
	return true
}

func (e *Error) IsNotError() bool {
	return !e.IsError()
}

func (e *Error) IsWarning() bool {
	if e.warning == nil {
		return false
	}
	return true
}

func (e *Error) IsNotWarning() bool {
	return !e.IsWarning()
}

func (e *Error) ExitIfError() {
	if e.IsError() {
		e.Print()
		os.Exit(1)
	}
}

func (e *Error) Print() {
	str := e.String()
	if str == "" {
		return
	}
	fmt.Print(str)
}

func (e Error) String() string {
	if e.err != nil {
		return fmt.Sprintf("%sERROR: %v\n", e.prefix, e.err)
	}
	if e.warning != nil {
		return fmt.Sprintf("%sWARNING: %v\n", e.prefix, e.warning)
	}
	return ""
}
