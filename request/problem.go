package request

type Problem struct {
	Title        string        `json:"title" binding:"required,max=20"`
	TimeLimit    int           `json:"time_limit"  binding:"required"`
	MemoryLimit  int           `json:"memory_limit"  binding:"required"`
	Description  string        `json:"description"`
	Input        string        `json:"input"`
	Output       string        `json:"output"`
	SampleInput  string        `json:"sample_input"`
	SampleOutput string        `json:"sample_output"`
	Spj          int           `json:"spj"`
	Hint         string        `json:"hint"`
	Level        int           `json:"level",gte=0, lte=2`
	Tags         []interface{} `json:"tags"`
}

type ProblemData struct {
	FileName string `json:"filename" binding:"required,max=20"`
}

type ProblemDataContent struct {
	Content string `json:"content" binding:"required"`
}
