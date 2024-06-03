package uerr

import "fmt"

type ErrorNotFound struct {
	Err error
}

func (e *ErrorNotFound) Error() string {
	return fmt.Errorf("error, not found: %w", e.Err).Error()
}

func (e *ErrorNotFound) Unwrap() error {
	return e.Err
}

func (e *ErrorNotFound) Is(target error) bool {
	_, ok := target.(*ErrorNotFound)
	return ok
}

type ErrorCannotInsert struct {
	Err error
}

func (e *ErrorCannotInsert) Error() string {
	return fmt.Errorf("error, cannot insert: %v", e.Err).Error()
}

func (e *ErrorCannotInsert) Unwrap() error {
	return e.Err
}

func (e *ErrorCannotInsert) Is(target error) bool {
	_, ok := target.(*ErrorCannotInsert)
	return ok
}

type ErrorCannotDelete struct {
	Err error
}

func (e *ErrorCannotDelete) Error() string {
	return fmt.Errorf("error, cannot delete: %w", e.Err).Error()
}

func (e *ErrorCannotDelete) Unwrap() error {
	return e.Err
}

func (e *ErrorCannotDelete) Is(target error) bool {
	_, ok := target.(*ErrorCannotDelete)
	return ok
}

type ErrorCannotUpdate struct {
	Err error
}

func (e *ErrorCannotUpdate) Error() string {
	return fmt.Errorf("error, cannot update: %w", e.Err).Error()
}

func (e *ErrorCannotUpdate) Unwrap() error {
	return e.Err
}

func (e *ErrorCannotUpdate) Is(target error) bool {
	_, ok := target.(*ErrorCannotUpdate)
	return ok
}
