package errorx

import (
	"errors"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// metadata key
const mdBizCode = "biz_code"

// ToGRPC converts an error to a gRPC status error with business code in details.
func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var ae *AppError
	if !errors.As(err, &ae) {
		// 非业务错误：对外统一 internal message，避免泄露内部信息
		s := status.New(codes.Internal, "服务器内部错误")
		st, _ := s.WithDetails(&errdetails.ErrorInfo{
			Reason: "INTERNAL_ERROR",
			Metadata: map[string]string{
				mdBizCode: strconv.Itoa(CodeInternalError),
			},
		})
		return st.Err()
	}

	grpcCode := grpcCodeFromBiz(ae.Code)

	s := status.New(grpcCode, ae.Msg)
	st, _ := s.WithDetails(&errdetails.ErrorInfo{
		Reason: "BIZ_ERROR",
		Metadata: map[string]string{
			mdBizCode: strconv.Itoa(ae.Code),
		},
	})
	return st.Err()
}

func grpcCodeFromBiz(bizCode int) codes.Code {
	switch bizCode {
	case CodeBadRequest:
		return codes.InvalidArgument
	case CodeUnauthorized:
		return codes.Unauthenticated
	case CodeForbidden:
		return codes.PermissionDenied
	case CodeNotFound:
		return codes.NotFound
	case CodeConflict:
		return codes.AlreadyExists
	case CodeTooManyRequests:
		return codes.ResourceExhausted
	case CodeUnavailable:
		return codes.Unavailable
	default:
		// 业务码（1000+）一般按语义映射：
		// - NotFound 类：NotFound
		// - AlreadyExists 类：AlreadyExists
		// - 其余默认 FailedPrecondition/Unknown 也行
		switch bizCode {
		case CodeUserNotFound, CodeCatNotFound, CodeTaskNotFound, CodeAdoptionNotFound, CodePostNotFound:
			return codes.NotFound
		case CodeUserAlreadyExists, CodeTaskAlreadyClaimed, CodeTaskFull, CodeAdoptionAlreadyApplied:
			return codes.AlreadyExists
		case CodePasswordWrong, CodeTokenInvalid, CodeTokenExpired:
			return codes.Unauthenticated
		case CodePermissionDenied, CodeTaskNotOwned, CodeInsufficientPoints:
			return codes.PermissionDenied
		default:
			return codes.Unknown
		}
	}
}

// FromGRPC extracts AppError from a gRPC error. If not present, returns internal error.
func FromGRPC(err error) *AppError {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return ErrInternalServer
	}

	// default
	msg := st.Message()
	code := CodeInternalError

	for _, d := range st.Details() {
		if ei, ok := d.(*errdetails.ErrorInfo); ok {
			if v, ok := ei.Metadata[mdBizCode]; ok {
				if n, e := strconv.Atoi(v); e == nil {
					code = n
				}
			}
		}
	}

	return New(code, msg)
}
