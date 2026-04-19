package kafka

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/luyb177/meow-nook/common/mail"
)

//go:embed templates/dlq_email.html
var dlqEmailTemplate string

var dlqTmpl = template.Must(
	template.New("dlq").Parse(dlqEmailTemplate),
)

type dlqEmailData struct {
	ServiceName string
	TaskID      string
	Retry       int
	MaxRetry    int
	Reason      string
	Stage       string
	Error       string
	CreatedAt   string
	NextRunAt   string
	Data        string
}

type EmailDLQNotifier struct {
	mailer      *mail.Mailer
	serviceName string
	to          []string
}

func NewEmailDLQNotifier(serviceName string, mailer *mail.Mailer, to []string) Notifier {
	return &EmailDLQNotifier{
		serviceName: serviceName,
		mailer:      mailer,
		to:          to,
	}
}

func (n *EmailDLQNotifier) Notify(ctx context.Context, env *Envelope) error {
	if len(n.to) == 0 {
		return nil // 没配收件人就不发，避免报错（同时靠日志）
	}

	subject := fmt.Sprintf("[DLQ][%s] task failed: %s", n.serviceName, env.TaskID)

	data := dlqEmailData{
		ServiceName: n.serviceName,
		TaskID:      env.TaskID,
		Retry:       env.Retry,
		MaxRetry:    env.MaxRetry,
		Reason:      env.DLQReason,
		Stage:       env.DLQStage,
		Error:       env.Error,
		CreatedAt:   formatUnix(env.CreatedAt),
		NextRunAt:   formatUnix(env.NextRunAt),
		Data:        truncate(string(env.Data), 8192),
	}

	var buf bytes.Buffer
	if err := dlqTmpl.Execute(&buf, data); err != nil {
		return err
	}

	body := buf.String()
	return n.mailer.Send(subject, body, n.to)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n...TRUNCATED..."
}

func formatUnix(ts int64) string {
	if ts <= 0 {
		return "-"
	}
	return time.Unix(ts, 0).Format(time.RFC3339)
}

func joinLines(ss []string) string {
	return strings.Join(ss, "\n")
}
