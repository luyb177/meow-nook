// Package errorx provides a unified error abstraction for the meow-nook backend.
// Business errors carry a machine-readable Code plus a human-readable Msg that
// can be sent back to the client, while underlying Go errors are wrapped so
// that the full context is still visible in logs.
package errorx

import (
	"errors"
	"fmt"
	"net/http"
)

// ──────────────────────────────────────────────
// Error codes
// ──────────────────────────────────────────────

const (
	// Success
	CodeOK = 0

	// Client-side errors (4xx)
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeConflict        = 409
	CodeUnprocessable   = 422
	CodeTooManyRequests = 429

	// Server-side errors (5xx)
	CodeInternalError  = 500
	CodeNotImplemented = 501
	CodeUnavailable    = 503

	// Business-domain codes (starting at 1000)
	CodeUserNotFound           = 1001
	CodeUserAlreadyExists      = 1002
	CodePasswordWrong          = 1003
	CodeTokenInvalid           = 1004
	CodeTokenExpired           = 1005
	CodeInsufficientPoints     = 1010
	CodeCatNotFound            = 1020
	CodeTaskNotFound           = 1030
	CodeTaskAlreadyClaimed     = 1031
	CodeTaskFull               = 1032
	CodeTaskNotOwned           = 1033
	CodeAdoptionNotFound       = 1040
	CodeAdoptionAlreadyApplied = 1041
	CodePostNotFound           = 1050
	CodePermissionDenied       = 1060
)

// codeToHTTPStatus maps business error codes to HTTP status codes.
var codeToHTTPStatus = map[int]int{
	// Success
	CodeOK: http.StatusOK,

	// Client-side errors (4xx)
	CodeBadRequest:      http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	CodeConflict:        http.StatusConflict,
	CodeUnprocessable:   http.StatusUnprocessableEntity,
	CodeTooManyRequests: http.StatusTooManyRequests,

	// Server-side errors (5xx)
	CodeInternalError:  http.StatusInternalServerError,
	CodeNotImplemented: http.StatusNotImplemented,
	CodeUnavailable:    http.StatusServiceUnavailable,

	// Business-domain codes (starting at 1000)
	CodeUserNotFound:           http.StatusNotFound,
	CodeUserAlreadyExists:      http.StatusConflict,
	CodePasswordWrong:          http.StatusUnauthorized,
	CodeTokenInvalid:           http.StatusUnauthorized,
	CodeTokenExpired:           http.StatusUnauthorized,
	CodeInsufficientPoints:     http.StatusForbidden,
	CodeCatNotFound:            http.StatusNotFound,
	CodeTaskNotFound:           http.StatusNotFound,
	CodeTaskAlreadyClaimed:     http.StatusConflict,
	CodeTaskFull:               http.StatusConflict,
	CodeTaskNotOwned:           http.StatusForbidden,
	CodeAdoptionNotFound:       http.StatusNotFound,
	CodeAdoptionAlreadyApplied: http.StatusConflict,
	CodePostNotFound:           http.StatusNotFound,
	CodePermissionDenied:       http.StatusForbidden,
}

// ──────────────────────────────────────────────
// AppError type
// ──────────────────────────────────────────────

// AppError is the standard error type used throughout the meow-nook backend.
// It carries a business error code, a user-facing message, and an optional
// cause for internal logging.
type AppError struct {
	Code  int
	Msg   string
	cause error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("code=%d msg=%s cause=%v", e.Code, e.Msg, e.cause)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// Unwrap enables errors.Is / errors.As traversal.
func (e *AppError) Unwrap() error {
	return e.cause
}

// HTTPStatus returns the HTTP status code that best represents this error.
func (e *AppError) HTTPStatus() int {
	if s, ok := codeToHTTPStatus[e.Code]; ok {
		return s
	}
	return http.StatusInternalServerError
}

// ──────────────────────────────────────────────
// Constructors
// ──────────────────────────────────────────────

// New creates a new AppError with the given code and message.
func New(code int, msg string) *AppError {
	return &AppError{Code: code, Msg: msg}
}

// Wrap creates an AppError that wraps an underlying error.
func Wrap(code int, msg string, cause error) *AppError {
	return &AppError{Code: code, Msg: msg, cause: cause}
}

// ──────────────────────────────────────────────
// Predefined sentinel errors
// ──────────────────────────────────────────────

var (
	// ─────────────────────────────
	// Common errors
	// ─────────────────────────────

	ErrBadRequest     = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized   = New(CodeUnauthorized, "请先登录")
	ErrForbidden      = New(CodeForbidden, "无权限访问")
	ErrNotFound       = New(CodeNotFound, "资源不存在")
	ErrInternalServer = New(CodeInternalError, "服务器内部错误")

	// ─────────────────────────────
	// Auth / Token errors
	// ─────────────────────────────

	ErrTokenInvalid = New(CodeTokenInvalid, "令牌无效")
	ErrTokenExpired = New(CodeTokenExpired, "令牌已过期")

	// ─────────────────────────────
	// User errors
	// ─────────────────────────────

	ErrUserNotFound      = New(CodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists = New(CodeUserAlreadyExists, "用户已存在")
	ErrPasswordWrong     = New(CodePasswordWrong, "密码错误")

	// ─────────────────────────────
	// Resource / Business errors
	// ─────────────────────────────

	ErrInsufficientPoints     = New(CodeInsufficientPoints, "积分不足")
	ErrCatNotFound            = New(CodeCatNotFound, "猫咪档案不存在")
	ErrTaskNotFound           = New(CodeTaskNotFound, "任务不存在")
	ErrTaskAlreadyClaimed     = New(CodeTaskAlreadyClaimed, "任务已被认领")
	ErrTaskFull               = New(CodeTaskFull, "任务认领人数已满")
	ErrTaskNotOwned           = New(CodeTaskNotOwned, "非本人认领的任务")
	ErrAdoptionNotFound       = New(CodeAdoptionNotFound, "领养申请不存在")
	ErrAdoptionAlreadyApplied = New(CodeAdoptionAlreadyApplied, "已申请过该猫咪的领养")
	ErrPostNotFound           = New(CodePostNotFound, "动态不存在")

	// ─────────────────────────────
	// Permission errors
	// ─────────────────────────────

	ErrPermissionDenied = New(CodePermissionDenied, "没有操作权限")

	// ─────────────────────────────
	// Server capability errors
	// ─────────────────────────────

	ErrNotImplemented = New(CodeNotImplemented, "功能暂未实现")
)

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return ErrInternalServer
}

// CodeOf returns the business code of err if it is an *AppError, or
// CodeInternalError otherwise.
func CodeOf(err error) int {
	var e *AppError
	if errors.As(err, &e) {
		return e.Code
	}
	return CodeInternalError
}

// MsgOf returns the user-facing message of err if it is an *AppError, or
// a generic message otherwise.
func MsgOf(err error) string {
	var e *AppError
	if errors.As(err, &e) {
		return e.Msg
	}
	return "服务器内部错误"
}
