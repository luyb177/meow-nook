package logic

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PBTimeToString 把 grpc timestamp 转成字符串，nil 安全
func PBTimeToString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format("2006-01-02 15:04:05")
}
