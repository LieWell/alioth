package core

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	SuccessCode    = "0"
	SuccessMessage = "success"
)

var SuccessResponse = NewWithSuccess(nil)

type IWrapperError interface {
	RawError() error
	StatusCode() int
	Code() string
	Message() string
	Data() any
}

type WrapperError struct {
	rawError   error  // 原始错误,用来还原现场
	statusCode int    // http 状态码
	C          string `json:"code"`    // 错误码
	M          string `json:"message"` // 错误原因
	D          any    `json:"data"`    // 返回值信息
}

func (err *WrapperError) RawError() error {
	if err.rawError == nil {
		return fmt.Errorf("buildin error: %v", err.M)
	}
	return err.rawError
}

func (err *WrapperError) Code() string {
	return err.C
}

func (err *WrapperError) Message() string {
	return err.M
}

func (err *WrapperError) StatusCode() int {
	return err.statusCode
}

func (err *WrapperError) Error() string {
	return err.M
}

func (err *WrapperError) Data() any {
	return err.D
}

var _ IWrapperError = (*WrapperError)(nil)
var _ error = (*WrapperError)(nil)

// NewWithSuccess 返回一个 "成功" 的错误
func NewWithSuccess(data any) *WrapperError {
	return &WrapperError{
		rawError:   nil,
		statusCode: http.StatusOK,
		C:          SuccessCode,
		D:          data,
		M:          SuccessMessage,
	}
}

// NewWithError 返回一个真正的错误
func NewWithError(raw error, statusCode int, code, message string) *WrapperError {
	return &WrapperError{
		rawError:   raw,
		statusCode: statusCode,
		C:          code,
		M:          message,
		D:          nil,
	}
}

func BadRequestError(code, message string) error {
	return NewWithError(nil, http.StatusBadRequest, code, message)
}

func SimpleBadRequestError() error {
	c := http.StatusBadRequest
	m := http.StatusText(c)
	return NewWithError(nil, c, strconv.Itoa(c), m)
}

func InternalServerError(code, message string) error {
	return NewWithError(nil, http.StatusInternalServerError, code, message)
}

func SimpleInternalServerError() error {
	c := http.StatusInternalServerError
	m := http.StatusText(c)
	return NewWithError(nil, c, strconv.Itoa(c), m)
}

func NotFoundError(code, message string) error {
	return NewWithError(nil, http.StatusNotFound, code, message)
}

func UnauthorizedError(code, message string) error {
	return NewWithError(nil, http.StatusUnauthorized, code, message)
}

func SimpleUnauthorizedError() error {
	c := http.StatusUnauthorized
	m := http.StatusText(c)
	return NewWithError(nil, c, strconv.Itoa(c), m)
}

func ForbiddenError(code, message string) error {
	return NewWithError(nil, http.StatusForbidden, code, message)
}

func ConflictError(code, message string) error {
	return NewWithError(nil, http.StatusConflict, code, message)
}
