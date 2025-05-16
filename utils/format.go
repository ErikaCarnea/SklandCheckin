package utils

import (
	"fmt"
	"github.com/HeathErika/Skland/models"
	"strconv"
	"time"
)

func FormatAttendanceResult(result *models.AttendanceResult) string {
	ts, _ := strconv.ParseInt(result.Data.Timestamp, 10, 64)
	localTime := time.Unix(ts, 0).Local().Format("2006-01-02 15:04:05")

	var awards string
	for _, award := range result.Data.Awards {
		awards += fmt.Sprintf("获得奖励：%s x%d (类型：%s)\n",
			award.Resource.Name,
			award.Count,
			award.Type)
	}

	return fmt.Sprintf("[%s] 签到成功！\n%s", localTime, awards)
}
