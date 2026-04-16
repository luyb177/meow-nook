package errorx_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/luyb177/meow-nook/common/errorx"
)

func TestNew(t *testing.T) {
	err := errorx.New(errorx.CodeBadRequest, "bad input")
	if err.Code != errorx.CodeBadRequest {
		t.Fatalf("expected code %d, got %d", errorx.CodeBadRequest, err.Code)
	}
	if err.Msg != "bad input" {
		t.Fatalf("expected msg %q, got %q", "bad input", err.Msg)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying cause")
	err := errorx.Wrap(errorx.CodeInternalError, "something went wrong", cause)
	if !errors.Is(err, cause) {
		t.Fatal("wrapped error should unwrap to the cause")
	}
	if err.Code != errorx.CodeInternalError {
		t.Fatalf("expected code %d, got %d", errorx.CodeInternalError, err.Code)
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		code       int
		wantStatus int
	}{
		{errorx.CodeBadRequest, http.StatusBadRequest},
		{errorx.CodeUnauthorized, http.StatusUnauthorized},
		{errorx.CodeForbidden, http.StatusForbidden},
		{errorx.CodeNotFound, http.StatusNotFound},
		{errorx.CodeInternalError, http.StatusInternalServerError},
		{errorx.CodeUserNotFound, http.StatusNotFound},
		{errorx.CodeTaskAlreadyClaimed, http.StatusConflict},
		{errorx.CodeInsufficientPoints, http.StatusForbidden},
	}
	for _, tc := range tests {
		err := errorx.New(tc.code, "test")
		if got := err.HTTPStatus(); got != tc.wantStatus {
			t.Errorf("code %d: want HTTP %d, got %d", tc.code, tc.wantStatus, got)
		}
	}
}

func TestIsAppError(t *testing.T) {
	appErr := errorx.New(errorx.CodeNotFound, "not found")
	if !errorx.IsAppError(appErr) {
		t.Fatal("expected IsAppError to return true for *AppError")
	}

	plain := errors.New("plain error")
	if errorx.IsAppError(plain) {
		t.Fatal("expected IsAppError to return false for plain error")
	}
}

func TestCodeOf(t *testing.T) {
	appErr := errorx.New(errorx.CodeCatNotFound, "cat not found")
	if got := errorx.CodeOf(appErr); got != errorx.CodeCatNotFound {
		t.Fatalf("want %d, got %d", errorx.CodeCatNotFound, got)
	}

	plain := errors.New("plain")
	if got := errorx.CodeOf(plain); got != errorx.CodeInternalError {
		t.Fatalf("want %d, got %d", errorx.CodeInternalError, got)
	}
}

func TestSentinelErrors(t *testing.T) {
	sentinels := []*errorx.AppError{
		errorx.ErrBadRequest,
		errorx.ErrUnauthorized,
		errorx.ErrForbidden,
		errorx.ErrNotFound,
		errorx.ErrInternalServer,
		errorx.ErrTokenInvalid,
		errorx.ErrTokenExpired,
		errorx.ErrUserNotFound,
		errorx.ErrUserAlreadyExists,
		errorx.ErrPasswordWrong,
		errorx.ErrInsufficientPoints,
		errorx.ErrCatNotFound,
		errorx.ErrTaskNotFound,
		errorx.ErrTaskAlreadyClaimed,
		errorx.ErrTaskFull,
		errorx.ErrTaskNotOwned,
		errorx.ErrAdoptionNotFound,
		errorx.ErrAdoptionAlreadyApplied,
		errorx.ErrPostNotFound,
		errorx.ErrPermissionDenied,
	}
	for _, s := range sentinels {
		if s.Msg == "" {
			t.Errorf("sentinel with code %d has empty Msg", s.Code)
		}
		if s.HTTPStatus() == 0 {
			t.Errorf("sentinel with code %d has zero HTTPStatus", s.Code)
		}
	}
}
