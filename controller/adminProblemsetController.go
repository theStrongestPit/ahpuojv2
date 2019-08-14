package controller

import (
	"ahpuoj/model"
	"ahpuoj/utils"
	"database/sql"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ImportProblemSet(c *gin.Context) {

	user, _ := GetUserInstance(c)

	filehead, err := c.FormFile("file")
	if utils.CheckError(c, err, "文件上传失败") != nil {
		return
	}
	file, err := filehead.Open()
	if utils.CheckError(c, err, "文件打开失败") != nil {
		return
	}
	fps, err := utils.ImportFps(file)
	if utils.CheckError(c, err, "问题导入失败") != nil {
		return
	}
	var infos []string

	for _, item := range fps.Item {
		timeLimit, _ := strconv.Atoi(item.TimeLimit.Content)
		memoryLimit, _ := strconv.Atoi(item.MemoryLimit.Content)
		if item.TimeLimit.Unit == "ms" {
			timeLimit /= 1000
		}
		if item.MemoryLimit.Unit == "kb" {
			memoryLimit /= 1024
		}
		problem := model.Problem{
			Title:        item.Title,
			Description:  sql.NullString{item.Description, true},
			Input:        sql.NullString{item.Input, true},
			Output:       sql.NullString{item.Output, true},
			SampleInput:  sql.NullString{item.SampleInput, true},
			SampleOutput: sql.NullString{item.SampleOutput, true},
			Hint:         sql.NullString{item.Hint, true},
			TimeLimit:    timeLimit,
			MemoryLimit:  memoryLimit,
		}
		err := problem.Save()

		if err != nil {
			infos = append(infos, "问题"+problem.Title+"导入失败")
		} else {
			infos = append(infos, "问题"+problem.Title+"导入成功")
			pid := problem.Id
			dataDir, _ := utils.GetCfg().GetValue("project", "datadir")
			baseDir := dataDir + "/" + strconv.Itoa(pid)
			err = os.MkdirAll(baseDir, 0777)
			if err != nil {
				utils.Consolelog(err.Error())
			}
			if len(item.SampleInput) > 0 {
				utils.Mkdata(pid, "sample.in", item.SampleInput)
			}
			if len(item.SampleOutput) > 0 {
				utils.Mkdata(pid, "sample.out", item.SampleOutput)
			}
			for index, testin := range item.TestInput {
				utils.Mkdata(pid, "test"+strconv.Itoa(index)+".in", testin)
			}
			for index, testout := range item.TestOutput {
				utils.Mkdata(pid, "test"+strconv.Itoa(index)+".out", testout)
			}
			// 提交默认答案
			for _, source := range item.Solution {

				var languageId int
				// 查找 language 的index
				for k, v := range utils.LanguageName {
					if v == source.Language {
						languageId = k
						break
					}
				}
				solution := model.Solution{
					ProblemId:  problem.Id,
					TeamId:     0,
					UserId:     user.Id,
					ContestId:  0,
					Num:        0,
					IP:         c.ClientIP(),
					Language:   languageId,
					CodeLength: len(source.Content),
				}
				err := solution.Save()
				if utils.CheckError(c, err, "保存提交记录失败") != nil {
					return
				}
				sourceCode := model.SourceCode{
					SolutionId: solution.Id,
					Source:     source.Content,
				}
				err = sourceCode.Save()
				if utils.CheckError(c, err, "保存代码记录失败") != nil {
					return
				}

				// 更新提交状态为等待评判
				_, err = DB.Exec("update solution set result = 0 where solution_id = ?", solution.Id)
				utils.Consolelog(err)

			}
		}
	}

	c.JSON(200, gin.H{
		"message": "操作成功",
		"info":    infos,
	})
}
