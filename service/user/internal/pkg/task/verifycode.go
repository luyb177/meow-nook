package task

const (
	TypeSendVerifyCode = "user.send_verify_code"
)

type VerifyCode struct {
	To   string
	Code string
}
