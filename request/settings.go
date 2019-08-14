package request

type Settings struct {
	EnableIssue bool `json:"enable_issue" binding:"required"`
}
