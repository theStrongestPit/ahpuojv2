package request

type Casbin struct {
	Rolename string `json:"rolename" binding:"required,max=40"`
	Path     string `json:"path" binding:"required,max=40"`
	Method   string `json:"method" binding:"required,max=40"`
}
