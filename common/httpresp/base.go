package httpresp

import (
	"context"
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type EmptyData struct {
}

func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	// v 可能是 error，也可能是正常返回 data
	if err, ok := v.(error); ok {
		ae := errorx.AsAppError(err)
		status := ae.HTTPStatus()
		httpx.WriteJsonCtx(ctx, w, status, BaseResponse{Code: ae.Code, Msg: ae.Msg, Data: &EmptyData{}})
		return
	}

	httpx.OkJsonCtx(ctx, w, BaseResponse{
		Code: errorx.CodeOK,
		Msg:  "ok",
		Data: v,
	})
}
