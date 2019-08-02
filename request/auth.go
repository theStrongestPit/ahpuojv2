package request

type Register struct {
	Email           string `json:"email" binding:"required,email,max=40"`
	Username        string `json:"username" binding:"required,ascii,max=20"`
	Nick            string `json:"nick" binding:"required,max=20"`
	Password        string `json:"password" binding:"required,ascii,min=6,max=20,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmpassword" binding:"required`
}
type Login struct {
	Username string `json:"username" binding:"required,ascii,max=20"`
	Password string `json:"password" binding:"required,ascii,min=6,max=20"`
}

type ResetPassByToken struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,ascii,min=6,max=20,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmpassword" binding:"required`
}
