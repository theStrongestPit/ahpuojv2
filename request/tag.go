package request

type Tag struct {
	Name string `json:"name" binding:"required,max=20"`
}
