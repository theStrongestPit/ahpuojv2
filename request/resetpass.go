package request

type ResetPass struct {
	Email string `json:"email" binding:"required,email,max=40"`
}
