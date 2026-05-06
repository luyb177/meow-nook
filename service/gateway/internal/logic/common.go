package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	httpmd "github.com/luyb177/meow-nook/common/middleware/http"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PBTimeToString 把 grpc timestamp 转成字符串，nil 安全
func PBTimeToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format(time.DateTime)
}

func ParseTimePtr(s string) *timestamppb.Timestamp {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t2, err2 := time.Parse(time.RFC3339, s)
		if err2 != nil {
			return nil
		}
		return timestamppb.New(t2)
	}
	return timestamppb.New(t)
}

// GetUserID todo g u f t
func GetUserID(ctx context.Context) (int64, error) {
	claims, ok := httpmd.ClaimsFromContext(ctx)
	if !ok || claims == nil {
		return 0, errorx.ErrUnauthorized
	}

	if claims.UserID == 0 {
		return 0, errorx.ErrUnauthorized
	}

	return claims.UserID, nil
}
