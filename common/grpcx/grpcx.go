// Package grpcx provides shared gRPC utilities for meow-nook services.
// It contains server interceptors and client-side helpers for converting
// between *errorx.AppError and gRPC status errors with ErrorDetail details.
package grpcx

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/errorx"
	errordetailv1 "github.com/luyb177/meow-nook/common/pb/errordetail/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AppErrorInterceptor is a gRPC unary server interceptor that converts
// *errorx.AppError values returned by handlers into gRPC status errors with
// an ErrorDetail message attached via WithDetails.  The ErrorDetail lets the
// HTTP gateway reconstruct the exact business {code, msg} without
// string-parsing the status message.
//
// If the error is not an *errorx.AppError the handler's error is mapped to
// codes.Internal; the original error is NOT forwarded to the client to avoid
// leaking internal details.
func AppErrorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	var ae *errorx.AppError
	if !errors.As(err, &ae) {
		// Non-AppError: do not expose internal details.
		return nil, status.Error(codes.Internal, "内部服务错误")
	}

	detail := &errordetailv1.ErrorDetail{
		Code: int32(ae.Code),
		Msg:  ae.Msg,
	}

	st, stErr := status.New(appCodeToGRPCCode(ae.Code), ae.Msg).WithDetails(detail)
	if stErr != nil {
		// Fallback: return status without details.
		return nil, status.Error(appCodeToGRPCCode(ae.Code), ae.Msg)
	}
	return nil, st.Err()
}

// ExtractAppError attempts to extract an *errorx.AppError from a gRPC status
// error by inspecting the status details for an ErrorDetail message.
//
// Decoding order:
//  1. If err is a gRPC status error with an ErrorDetail detail, reconstruct
//     an *errorx.AppError from it.
//  2. If the status has no ErrorDetail, build an *errorx.AppError from the
//     gRPC code and status message as a fallback.
//  3. If err is not a gRPC status error at all, return nil (caller should
//     treat as internal error).
func ExtractAppError(err error) *errorx.AppError {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return nil
	}
	if st.Code() == codes.OK {
		return nil
	}

	// Walk the status details looking for our ErrorDetail message.
	for _, d := range st.Details() {
		if ed, ok := d.(*errordetailv1.ErrorDetail); ok {
			return errorx.New(int(ed.Code), ed.Msg)
		}
	}

	// Fallback: no ErrorDetail found – derive AppError from the gRPC status.
	return errorx.New(grpcCodeToAppCode(st.Code()), st.Message())
}

// appCodeToGRPCCode maps an errorx business code to the most appropriate
// gRPC status code.
func appCodeToGRPCCode(code int) codes.Code {
	switch code {
	case errorx.CodeBadRequest, errorx.CodeUnprocessable:
		return codes.InvalidArgument
	case errorx.CodeUnauthorized,
		errorx.CodePasswordWrong,
		errorx.CodeTokenInvalid,
		errorx.CodeTokenExpired:
		return codes.Unauthenticated
	case errorx.CodeForbidden,
		errorx.CodeInsufficientPoints,
		errorx.CodePermissionDenied,
		errorx.CodeTaskNotOwned:
		return codes.PermissionDenied
	case errorx.CodeNotFound,
		errorx.CodeUserNotFound,
		errorx.CodeCatNotFound,
		errorx.CodeTaskNotFound,
		errorx.CodeAdoptionNotFound,
		errorx.CodePostNotFound:
		return codes.NotFound
	case errorx.CodeConflict,
		errorx.CodeUserAlreadyExists,
		errorx.CodeTaskAlreadyClaimed,
		errorx.CodeTaskFull,
		errorx.CodeAdoptionAlreadyApplied:
		return codes.AlreadyExists
	case errorx.CodeTooManyRequests:
		return codes.ResourceExhausted
	case errorx.CodeNotImplemented:
		return codes.Unimplemented
	case errorx.CodeUnavailable:
		return codes.Unavailable
	default:
		return codes.Internal
	}
}

// grpcCodeToAppCode maps a gRPC status code to an errorx business code.
// Used as a fallback when no ErrorDetail is present in the status.
func grpcCodeToAppCode(c codes.Code) int {
	switch c {
	case codes.InvalidArgument:
		return errorx.CodeBadRequest
	case codes.Unauthenticated:
		return errorx.CodeUnauthorized
	case codes.PermissionDenied:
		return errorx.CodeForbidden
	case codes.NotFound:
		return errorx.CodeNotFound
	case codes.AlreadyExists:
		return errorx.CodeConflict
	case codes.ResourceExhausted:
		return errorx.CodeTooManyRequests
	case codes.Unimplemented:
		return errorx.CodeNotImplemented
	case codes.Unavailable:
		return errorx.CodeUnavailable
	default:
		return errorx.CodeInternalError
	}
}
