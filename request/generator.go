package request

type CompeteAccount struct {
	Prefix string `json:"prefix" binding:"required,max=15"`
	Number int    `json:"number" binding:"required,min=1,max=100"`
}

type UserAccount struct {
	UserList string `json:"userlist" binding:"required"`
}
