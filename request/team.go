package request

type Team struct {
	Name string `json:"name" binding:"required,max=20"`
}

type TeamUsers struct {
	UserList string `json:"userlist" binding:"required"`
}
