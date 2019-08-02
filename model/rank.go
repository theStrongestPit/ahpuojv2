package model

import "time"

type RankItem struct {
	ProblemId int    `db:"problem_id"`
	TeamId    int    `db:"team_id"`
	UserId    int    `db:"user_id"`
	ContestId int    `db:"contest_id"`
	Num       int    `db:"num"`
	InDate    string `db:"in_date"`
	Result    int    `db:"result"`
	// 附加信息
	Username   string `db:"username"`
	Nick       string `db:"nick"`
	UserRole   string `db:"name"`
	UserAvatar string `db:"avatar"`
}

type UserRankInfo struct {
	Solved  int   `json:"solved"`
	Time    int   `json:"time"`
	WaCount []int `json:"wa_count"`
	AcTime  []int `json:"ac_time"`
	TeamId  int   `json:"team_id"`
	User    struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Nick     string `json:"nick"`
	} `json:"user"`
}

type TeamRankInfo struct {
	Solved  int   `json:"solved"`
	Time    int   `json:"time"`
	WaCount []int `json:"wa_count"`
	AcCount []int `json:"ac_count"`
	AcTime  []int `json:"ac_time"`
	Team    struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
}

type UserRankInfoList []UserRankInfo
type UserRankSortByTeam []UserRankInfo

type TeamRankInfoList []TeamRankInfo

func (userRankInfo *UserRankInfo) Add(rankItem RankItem, startTimeStr string) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, time.Local)
	if userRankInfo.AcTime[rankItem.Num-1] > 0 {
		return
	}
	if rankItem.Result != 4 {
		userRankInfo.WaCount[rankItem.Num-1]++
	} else {
		acTime, _ := time.ParseInLocation("2006-01-02 15:04:05", rankItem.InDate, time.Local)
		useTime := int(acTime.Unix() - startTime.Unix())
		userRankInfo.AcTime[rankItem.Num-1] = useTime
		userRankInfo.Solved++
		userRankInfo.Time += useTime + 1200*userRankInfo.WaCount[rankItem.Num-1]
	}
}

func (teamRankInfo *TeamRankInfo) Add(userRankInfo UserRankInfo) {
	teamRankInfo.Solved += userRankInfo.Solved
	teamRankInfo.Time += userRankInfo.Time
	for k, v := range userRankInfo.AcTime {
		if v > 0 {
			teamRankInfo.AcCount[k]++
			teamRankInfo.AcTime[k] += v
		}
	}

	for k, v := range userRankInfo.WaCount {
		if v > 0 {
			teamRankInfo.WaCount[k] += v
		}
	}
}

// 个人排名排序
func (uril UserRankInfoList) Len() int {
	return len(uril)
}

func (uril UserRankInfoList) Swap(i, j int) {
	uril[i], uril[j] = uril[j], uril[i]
}

func (uril UserRankInfoList) Less(i, j int) bool {
	if uril[i].Solved != uril[j].Solved {
		return uril[i].Solved > uril[j].Solved
	} else {
		return uril[i].Time < uril[j].Time
	}
}

// 按照teamid排序
func (uril UserRankSortByTeam) Len() int {
	return len(uril)
}

func (uril UserRankSortByTeam) Swap(i, j int) {
	uril[i], uril[j] = uril[j], uril[i]
}

func (uril UserRankSortByTeam) Less(i, j int) bool {
	return uril[i].TeamId < uril[j].TeamId
}

// 团队排名排序
func (tril TeamRankInfoList) Len() int {
	return len(tril)
}

func (tril TeamRankInfoList) Swap(i, j int) {
	tril[i], tril[j] = tril[j], tril[i]
}

func (tril TeamRankInfoList) Less(i, j int) bool {
	if tril[i].Solved != tril[j].Solved {
		return tril[i].Solved > tril[j].Solved
	} else {
		return tril[i].Time < tril[j].Time
	}
}
