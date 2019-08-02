package request

type Series struct {
	Name        string `json:"name" binding:"required,max=20"`
	Description string `json:"description"`
	TeamMode    int    `json:"team_mode" binding:"gte=0,lte=1"`
}
