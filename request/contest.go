package request

// 限定了团队模式下只能为私有，即不存在private为0同时team_mode为1，即team_mode<=private
type Contest struct {
	Name        string `json:"name" binding:"required,max=20"`
	StartTime   string `json:"start_time"  binding:"required"`
	EndTime     string `json:"end_time"  binding:"required"`
	Description string `json:"description"`
	Problems    string `json:"problems"`
	LangMask    int    `json:"langmask"`
	Private     int    `json:"private"  binding:"gte=0,lte=1"`
	TeamMode    int    `json:"team_mode" ltefield=Private`
}

type ContestUsers struct {
	UserList string `json:"userlist" binding:"required"`
}
