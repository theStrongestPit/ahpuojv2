package model

import (
	"ahpuoj/utils"
	"database/sql"
)

type Issue struct {
	Id        int    `db:"id"`
	Title     string `db:"title"`
	ProblemId int    `db:"problem_id"`
	UserId    int    `db:"user_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	IsDeleted int    `db:"is_deleted"`
	// 附加信息
	Username     string         `db:"username"`
	Nick         string         `db:"nick"`
	UserAvatar   string         `db:"avatar"`
	ReplyCount   int            `db:"reply_count"`
	ProblemTitle sql.NullString `db:"ptitle"`
}

type Reply struct {
	Id          int    `db:"id"`
	IssueId     int    `db:"issue_id"`
	UserId      int    `db:"user_id"`
	ReplyId     int    `db:"reply_id"`
	ReplyUserId int    `db:"reply_user_id"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	Content     string `db:"content"`
	IsDeleted   int    `db:"is_deleted"`
	Status      int    `db:"status"`
	// 附加信息
	Username      string `db:"username"`
	ReplyUserNick string `db:"rnick"`
	Nick          string `db:"nick"`
	UserAvatar    string `db:"avatar"`
	ReplyCount    int    `db:"reply_count"`
	SubReplys     []map[string]interface{}
}

func (issue *Issue) Save() error {
	result, err := DB.Exec(`insert into issue
	(title,problem_id,user_id,created_at,updated_at) 
	values (?,?,?,NOW(),NOW())`, issue.Title, issue.ProblemId, issue.UserId)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	issue.Id = utils.Int64to32(lastInsertId)
	return err
}

func (reply *Reply) Save() error {
	result, err := DB.Exec(`insert into reply
	(user_id,issue_id,reply_id,reply_user_id,content,created_at,updated_at) 
	values (?,?,?,?,?,NOW(),NOW())`, reply.UserId, reply.IssueId, reply.ReplyId, reply.ReplyUserId, reply.Content)
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	reply.Id = utils.Int64to32(lastInsertId)
	// 更新主题的最后更新时间
	DB.Exec("update issue set updated_at = NOW() where id = ?", reply.IssueId)
	return err
}

func (issue *Issue) Response() map[string]interface{} {
	var problemTitle string
	if issue.ProblemTitle.Valid {
		problemTitle = issue.ProblemTitle.String
	} else {
		problemTitle = ""
	}
	return map[string]interface{}{
		"id":          issue.Id,
		"title":       issue.Title,
		"reply_count": issue.ReplyCount,
		"updated_at":  issue.UpdatedAt,
		"user": map[string]interface{}{
			"id":       issue.UserId,
			"username": issue.Username,
			"nick":     issue.Nick,
			"avatar":   issue.UserAvatar,
		},
		"problem": map[string]interface{}{
			"id":    issue.ProblemId,
			"title": problemTitle,
		},
	}
}

func (reply *Reply) Response() map[string]interface{} {
	// 需要将图片地址转换为绝对地址
	reply.Content = utils.ConvertTextImgUrl(reply.Content)
	return map[string]interface{}{
		"id":              reply.Id,
		"content":         reply.Content,
		"issue_id":        reply.IssueId,
		"reply_id":        reply.ReplyId,
		"reply_user_id":   reply.ReplyUserId,
		"reply_user_nick": reply.ReplyUserNick,
		"reply_count":     reply.ReplyCount,
		"updated_at":      reply.UpdatedAt,
		"sub_replys":      reply.SubReplys,
		"user": map[string]interface{}{
			"id":       reply.UserId,
			"username": reply.Username,
			"nick":     reply.Nick,
			"avatar":   reply.UserAvatar,
		},
	}
}
