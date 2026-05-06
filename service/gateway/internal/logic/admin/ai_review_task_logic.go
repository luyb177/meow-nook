package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type AIReviewTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAIReviewTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIReviewTaskLogic {
	return &AIReviewTaskLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AIReviewTaskLogic) AiReviewTask(req *types.AiReviewTaskReq) (*types.AiReviewTaskResp, error) {
	cfg := l.svcCtx.Config.AIService
	baseURL := cfg.ModelBaseURL
	if baseURL == "" {
		return nil, errorx.New(errorx.CodeInternalError, "MODEL_BASE_URL 未配置")
	}
	if req.TaskId == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "taskId 不能为空", errorx.ErrBadRequest)
	}

	payload := map[string]any{
		"model":   cfg.ModelName,
		"task_id": req.TaskId,
		"content": req.Content,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, errorx.WrapInternal("序列化请求失败", err)
	}

	httpReq, err := http.NewRequestWithContext(l.ctx, http.MethodPost, baseURL, bytes.NewReader(b))
	if err != nil {
		return nil, errorx.WrapInternal("创建请求失败", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if cfg.ModelAPIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+cfg.ModelAPIKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, errorx.WrapInternal("调用 AI 审核失败", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errorx.WrapInternal("读取 AI 响应失败", err)
	}

	out := &types.AiReviewTaskResp{
		StatusCode: int32(httpResp.StatusCode),
	}

	var anyResp any
	if err := json.Unmarshal(raw, &anyResp); err == nil {
		out.Result = anyResp
	} else {
		out.Raw = string(raw)
	}

	if httpResp.StatusCode >= 400 {
		return nil, errorx.New(errorx.CodeInternalError, "AI 审核接口返回错误")
	}

	return out, nil
}
