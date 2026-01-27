package utils

import (
	"fmt"
	"strings"

	"github.com/ErikaCarnea/Skland/models"
)

func FormatArknightSignResult(result *models.AttendanceResult) string {
	var awards strings.Builder
	for _, award := range result.Data.Awards {
		fmt.Fprintf(&awards, "获得奖励：%s x%d (类型：%s)\n",
			award.Resource.Name,
			award.Count,
			award.Type)
	}

	return fmt.Sprintf("签到成功！\n%s", awards.String())
}

func FormatEndfieldSignResult(result *models.EndfieldResult) string {
	var awards strings.Builder
	if result.Code == 10001 {
		return "今日已签到，明日再来吧！"
	}
	for _, award := range result.Data.AwardIds {
		fmt.Fprintf(&awards, "获得奖励：%s x%d (类型：%d)\n",
			result.Data.ResourceInfoMap[award.Id].Name,
			result.Data.ResourceInfoMap[award.Id].Count,
			award.Type)
	}

	return fmt.Sprintf("签到成功！\n%s", awards.String())
}
