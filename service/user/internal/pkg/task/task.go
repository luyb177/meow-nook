// Package task defines Kafka task types for the user service.
// TaskID format: {type}:{biz}:{uuid}  (built by kafka.BuildTaskID).

package task

const (
	TypeSuccessTest = "user.success_test"
	TypeFailTest    = "user.fail_test"
)

type TestPayload struct {
	Name string `json:"name"`
}
