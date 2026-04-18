package grpcx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/grpcx"
	errordetailv1 "github.com/luyb177/meow-nook/common/pb/errordetail/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ──────────────────────────────────────────────
// AppErrorInterceptor tests
// ──────────────────────────────────────────────

func TestAppErrorInterceptor_NoError(t *testing.T) {
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return "ok", nil
	}
	resp, err := grpcx.AppErrorInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("unexpected response %v", resp)
	}
}

func TestAppErrorInterceptor_AppError(t *testing.T) {
	appErr := errorx.New(errorx.CodeUserNotFound, "用户不存在")
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, appErr
	}

	_, err := grpcx.AppErrorInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T", err)
	}

	// gRPC code should map to NotFound.
	if st.Code() != codes.NotFound {
		t.Errorf("expected codes.NotFound, got %v", st.Code())
	}

	// Status message should carry the human-readable msg.
	if st.Message() != "用户不存在" {
		t.Errorf("unexpected status message %q", st.Message())
	}

	// Details should contain our ErrorDetail.
	details := st.Details()
	if len(details) == 0 {
		t.Fatal("expected at least one detail in status")
	}
	ed, ok := details[0].(*errordetailv1.ErrorDetail)
	if !ok {
		t.Fatalf("expected *ErrorDetail in details, got %T", details[0])
	}
	if int(ed.Code) != errorx.CodeUserNotFound {
		t.Errorf("expected code %d, got %d", errorx.CodeUserNotFound, ed.Code)
	}
	if ed.Msg != "用户不存在" {
		t.Errorf("expected msg %q, got %q", "用户不存在", ed.Msg)
	}
}

func TestAppErrorInterceptor_NonAppError(t *testing.T) {
	plainErr := errors.New("database connection failed")
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, plainErr
	}

	_, err := grpcx.AppErrorInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error")
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected codes.Internal, got %v", st.Code())
	}
	// Must NOT forward the original plain error message to prevent info leakage.
	if st.Message() == plainErr.Error() {
		t.Error("interceptor must not forward raw internal error messages to clients")
	}
}

// ──────────────────────────────────────────────
// ExtractAppError tests
// ──────────────────────────────────────────────

func TestExtractAppError_WithDetails(t *testing.T) {
	detail := &errordetailv1.ErrorDetail{Code: int32(errorx.CodeCatNotFound), Msg: "猫咪档案不存在"}
	st, err := status.New(codes.NotFound, "猫咪档案不存在").WithDetails(detail)
	if err != nil {
		t.Fatalf("WithDetails: %v", err)
	}

	ae := grpcx.ExtractAppError(st.Err())
	if ae == nil {
		t.Fatal("expected *AppError, got nil")
	}
	if ae.Code != errorx.CodeCatNotFound {
		t.Errorf("expected code %d, got %d", errorx.CodeCatNotFound, ae.Code)
	}
	if ae.Msg != "猫咪档案不存在" {
		t.Errorf("expected msg %q, got %q", "猫咪档案不存在", ae.Msg)
	}
}

func TestExtractAppError_FallbackNoDetails(t *testing.T) {
	st := status.New(codes.PermissionDenied, "no permission")
	ae := grpcx.ExtractAppError(st.Err())
	if ae == nil {
		t.Fatal("expected *AppError, got nil")
	}
	// Fallback mapping: PermissionDenied → CodeForbidden.
	if ae.Code != errorx.CodeForbidden {
		t.Errorf("expected code %d, got %d", errorx.CodeForbidden, ae.Code)
	}
	if ae.Msg != "no permission" {
		t.Errorf("expected msg %q, got %q", "no permission", ae.Msg)
	}
}

func TestExtractAppError_NonStatusError(t *testing.T) {
	ae := grpcx.ExtractAppError(errors.New("plain error"))
	if ae != nil {
		t.Errorf("expected nil for non-status error, got %+v", ae)
	}
}

func TestExtractAppError_NilError(t *testing.T) {
	ae := grpcx.ExtractAppError(nil)
	if ae != nil {
		t.Errorf("expected nil for nil error, got %+v", ae)
	}
}

// ──────────────────────────────────────────────
// Round-trip: interceptor → ExtractAppError
// ──────────────────────────────────────────────

func TestRoundTrip(t *testing.T) {
	orig := errorx.New(errorx.CodeTaskAlreadyClaimed, "任务已被认领")
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, orig
	}

	_, grpcErr := grpcx.AppErrorInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if grpcErr == nil {
		t.Fatal("expected error from interceptor")
	}

	ae := grpcx.ExtractAppError(grpcErr)
	if ae == nil {
		t.Fatal("ExtractAppError returned nil")
	}
	if ae.Code != orig.Code {
		t.Errorf("code mismatch: want %d, got %d", orig.Code, ae.Code)
	}
	if ae.Msg != orig.Msg {
		t.Errorf("msg mismatch: want %q, got %q", orig.Msg, ae.Msg)
	}
}
