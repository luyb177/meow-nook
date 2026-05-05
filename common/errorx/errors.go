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
// Format: AABBCC
//
// AA: 一级分类（系统级 / 业务级）
// BB: 子模块
// CC: 具体错误编号
//
// 示例：
// 200101
// 20 = 用户模块
// 01 = 用户基础信息
// 01 = 用户不存在
// ──────────────────────────────────────────────

const (
	// Success
	CodeOK = 0

	// =========================================================
	// Common HTTP-style errors（兼容基础状态码）
	// =========================================================

	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeConflict        = 409
	CodeUnprocessable   = 422
	CodeTooManyRequests = 429

	// =========================================================
	// System Errors（10 ~ 19）
	// =========================================================

	// 10 - 通用系统错误
	CodeInternalError  = 100001 // 系统内部错误
	CodeNotImplemented = 100002 // 功能暂未实现
	CodeUnavailable    = 100003 // 服务暂不可用

	// 11 - 数据库错误
	CodeDatabaseQueryFailed  = 110101 // 数据查询失败
	CodeDatabaseInsertFailed = 110102 // 数据插入失败
	CodeDatabaseUpdateFailed = 110103 // 数据更新失败
	CodeDatabaseDeleteFailed = 110104 // 数据删除失败
	CodeDatabaseTxFailed     = 110105 // 数据事务失败

	// 12 - Redis错误
	CodeRedisGetFailed = 120101 // Redis读取失败
	CodeRedisSetFailed = 120102 // Redis写入失败

	// 14 - RPC调用错误
	CodeRPCFailed = 140101 // RPC调用失败

	// 18 - 权限认证错误
	CodeTokenInvalid = 180101 // Token无效
	CodeTokenExpired = 180102 // Token过期

	// =========================================================
	// Business Errors（20 ~ 29）
	// =========================================================

	// 20 - 用户模块
	CodeUserNotFound      = 200101 // 用户不存在
	CodeUserAlreadyExists = 200102 // 用户已存在
	CodePasswordWrong     = 200103 // 密码错误

	// 21 - 猫咪模块
	CodeCatNotFound = 210101 // 猫咪档案不存在

	// 22 - 任务模块
	CodeTaskNotFound       = 220101 // 任务不存在
	CodeTaskAlreadyClaimed = 220102 // 任务已被认领
	CodeTaskFull           = 220103 // 任务人数已满
	CodeTaskNotOwned       = 220104 // 非本人认领任务

	// 23 - 领养模块
	CodeAdoptionNotFound       = 230101 // 领养申请不存在
	CodeAdoptionAlreadyApplied = 230102 // 已申请过该猫咪领养

	// 24 - 动态模块
	CodePostNotFound = 240101 // 动态不存在

	// 26 - 积分模块
	CodeInsufficientPoints = 260101 // 积分不足

	// 27 - 权限模块
	CodePermissionDenied = 270101 // 权限不足
)

// codeToHTTPStatus maps business error codes to HTTP status codes.
var codeToHTTPStatus = map[int]int{
	// Success
	CodeOK: http.StatusOK,

	// HTTP compatibility
	CodeBadRequest:      http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	CodeConflict:        http.StatusConflict,
	CodeUnprocessable:   http.StatusUnprocessableEntity,
	CodeTooManyRequests: http.StatusTooManyRequests,

	// System errors
	CodeInternalError:        http.StatusInternalServerError,
	CodeNotImplemented:       http.StatusNotImplemented,
	CodeUnavailable:          http.StatusServiceUnavailable,
	CodeDatabaseQueryFailed:  http.StatusInternalServerError,
	CodeDatabaseInsertFailed: http.StatusInternalServerError,
	CodeDatabaseUpdateFailed: http.StatusInternalServerError,
	CodeDatabaseDeleteFailed: http.StatusInternalServerError,
	CodeDatabaseTxFailed:     http.StatusInternalServerError,
	CodeRedisGetFailed:       http.StatusInternalServerError,
	CodeRedisSetFailed:       http.StatusInternalServerError,
	CodeRPCFailed:            http.StatusInternalServerError,
	CodeTokenInvalid:         http.StatusUnauthorized,
	CodeTokenExpired:         http.StatusUnauthorized,

	// Business errors
	CodeUserNotFound:           http.StatusNotFound,
	CodeUserAlreadyExists:      http.StatusConflict,
	CodePasswordWrong:          http.StatusUnauthorized,
	CodeCatNotFound:            http.StatusNotFound,
	CodeTaskNotFound:           http.StatusNotFound,
	CodeTaskAlreadyClaimed:     http.StatusConflict,
	CodeTaskFull:               http.StatusConflict,
	CodeTaskNotOwned:           http.StatusForbidden,
	CodeAdoptionNotFound:       http.StatusNotFound,
	CodeAdoptionAlreadyApplied: http.StatusConflict,
	CodePostNotFound:           http.StatusNotFound,
	CodeInsufficientPoints:     http.StatusForbidden,
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
	// =========================================================
	// Common errors
	// =========================================================

	ErrBadRequest     = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized   = New(CodeUnauthorized, "请先登录")
	ErrForbidden      = New(CodeForbidden, "无权限访问")
	ErrNotFound       = New(CodeNotFound, "资源不存在")
	ErrInternalServer = New(CodeInternalError, "服务器内部错误")
	ErrUnavailable    = New(CodeUnavailable, "服务暂不可用")
	ErrNotImplemented = New(CodeNotImplemented, "功能暂未实现")

	// =========================================================
	// Auth / Token errors
	// =========================================================

	ErrTokenInvalid = New(CodeTokenInvalid, "令牌无效")
	ErrTokenExpired = New(CodeTokenExpired, "令牌已过期")

	// =========================================================
	// Database errors
	// =========================================================

	ErrDatabaseQuery  = New(CodeDatabaseQueryFailed, "数据库查询失败")
	ErrDatabaseInsert = New(CodeDatabaseInsertFailed, "数据库新增失败")
	ErrDatabaseUpdate = New(CodeDatabaseUpdateFailed, "数据库更新失败")
	ErrDatabaseDelete = New(CodeDatabaseDeleteFailed, "数据库删除失败")
	ErrDatabaseTx     = New(CodeDatabaseTxFailed, "数据库事务失败")

	// =========================================================
	// Redis errors
	// =========================================================

	ErrRedisGet = New(CodeRedisGetFailed, "缓存读取失败")
	ErrRedisSet = New(CodeRedisSetFailed, "缓存写入失败")

	// =========================================================
	// RPC errors
	// =========================================================

	ErrRPCFailed = New(CodeRPCFailed, "服务调用失败")

	// =========================================================
	// User errors
	// =========================================================

	ErrUserNotFound      = New(CodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists = New(CodeUserAlreadyExists, "用户已存在")
	ErrPasswordWrong     = New(CodePasswordWrong, "密码错误")

	// =========================================================
	// Cat errors
	// =========================================================

	ErrCatNotFound = New(CodeCatNotFound, "猫咪档案不存在")

	// =========================================================
	// Task errors
	// =========================================================

	ErrTaskNotFound       = New(CodeTaskNotFound, "任务不存在")
	ErrTaskAlreadyClaimed = New(CodeTaskAlreadyClaimed, "任务已被认领")
	ErrTaskFull           = New(CodeTaskFull, "任务认领人数已满")
	ErrTaskNotOwned       = New(CodeTaskNotOwned, "非本人认领的任务")

	// =========================================================
	// Adoption errors
	// =========================================================

	ErrAdoptionNotFound       = New(CodeAdoptionNotFound, "领养申请不存在")
	ErrAdoptionAlreadyApplied = New(CodeAdoptionAlreadyApplied, "已申请过该猫咪的领养")

	// =========================================================
	// Post errors
	// =========================================================

	ErrPostNotFound = New(CodePostNotFound, "动态不存在")

	// =========================================================
	// Points errors
	// =========================================================

	ErrInsufficientPoints = New(CodeInsufficientPoints, "积分不足")

	// =========================================================
	// Permission errors
	// =========================================================

	ErrPermissionDenied = New(CodePermissionDenied, "没有操作权限")
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

func WrapInternal(msg string, err error) *AppError {
	return Wrap(CodeInternalError, msg, err)
}

func WrapDBQuery(msg string, err error) *AppError {
	return Wrap(CodeDatabaseQueryFailed, msg, err)
}

func WrapDBInsert(msg string, err error) *AppError {
	return Wrap(CodeDatabaseInsertFailed, msg, err)
}

func WrapDBUpdate(msg string, err error) *AppError {
	return Wrap(CodeDatabaseUpdateFailed, msg, err)
}

func WrapDBDelete(msg string, err error) *AppError {
	return Wrap(CodeDatabaseDeleteFailed, msg, err)
}

func WrapDBTx(msg string, err error) *AppError {
	return Wrap(CodeDatabaseTxFailed, msg, err)
}

func WrapRedisGet(msg string, err error) *AppError {
	return Wrap(CodeRedisGetFailed, msg, err)
}

func WrapRedisSet(msg string, err error) *AppError {
	return Wrap(CodeRedisSetFailed, msg, err)
}

func WrapRPC(msg string, err error) *AppError {
	return Wrap(CodeRPCFailed, msg, err)
}
