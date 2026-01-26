package utils

import (
	"fmt"

	"github.com/ErikaCarnea/Skland/models"
)

func FormatArknightSignResult(result *models.AttendanceResult) string {
	var awards string
	for _, award := range result.Data.Awards {
		awards += fmt.Sprintf("获得奖励：%s x%d (类型：%s)\n",
			award.Resource.Name,
			award.Count,
			award.Type)
	}

	return fmt.Sprintf("签到成功！\n%s", awards)
}

func FormatEndfieldSignResult(result *models.EndfieldResult) string {
	var awards string
	for _, award := range result.Data.AwardIds {
		awards += fmt.Sprintf("获得奖励：%s x%d (类型：%d)\n",
			result.Data.ResourceInfoMap[award.Id].Name,
			result.Data.ResourceInfoMap[award.Id].Count,
			award.Type)
	}

	return fmt.Sprintf("签到成功！\n%s", awards)
}
