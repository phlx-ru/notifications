// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsInternalError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_INTERNAL_ERROR.String() && e.Code == 500
}

func ErrorInternalError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_INTERNAL_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsInvalidRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_INVALID_REQUEST.String() && e.Code == 400
}

func ErrorInvalidRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_INVALID_REQUEST.String(), fmt.Sprintf(format, args...))
}

func IsNotificationNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NOTIFICATION_NOT_FOUND.String() && e.Code == 404
}

func ErrorNotificationNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, ErrorReason_NOTIFICATION_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}
