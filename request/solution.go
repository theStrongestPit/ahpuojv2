package request

type Solution struct {
	ProblemId int    `json:"problem_id" binding:"required"`
	Language  int    `json:"language" binding:"gte=0,lte=17"`
	ContestId int    `json:"contest_id"`
	Num       int    `json:"num" binding:"omitempty,gte=0"`
	Source    string `json:"source"  binding:"required,min=2,max=65535"`
}

type TestRun struct {
	Language  int    `json:"language" binding:"gte=0,lte=17"`
	InputText string `json:"input_text"  binding:"max=65535"`
	Source    string `json:"source"  binding:"required,min=2,max=65535"`
}
