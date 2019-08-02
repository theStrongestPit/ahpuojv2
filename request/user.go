package request

type UserPass struct {
	Password string `json:"password" binding:"required,ascii,min=6,max=20"`
}

type UserNick struct {
	Nick string `json:"nick" binding:"required,max=20"`
}

type UserResetPassword struct {
	OldPassword     string `json:"oldpassword" binding:"required,ascii,min=6,max=20"`
	Password        string `json:"password" binding:"required,ascii,min=6,max=20"`
	ConfirmPassword string `json:"confirmpassword" binding:"required,ascii,min=6,max=20"`
}
