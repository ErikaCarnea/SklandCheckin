package models

// GameID 定义森空岛支持的板块/游戏 ID。
type GameID int

const (
	GameArknights GameID = 1   // 明日方舟
	GameFromStars GameID = 2   // 来自星辰
	GameEndfield  GameID = 3   // 明日方舟: 终末地
	GamePopucom   GameID = 4   // 泡姆泡姆
	GameNastPort  GameID = 100 // 纳斯特港
	GameCoreBlaze GameID = 101 // 开拓芯
)

// String 返回游戏中文名。
func (g GameID) String() string {
	if name, ok := GameIDNames[g]; ok {
		return name
	}
	return "未知游戏"
}

// GameIDNames 游戏 ID → 中文名映射。
var GameIDNames = map[GameID]string{
	GameArknights: "明日方舟",
	GameFromStars: "来自星辰",
	GameEndfield:  "明日方舟: 终末地",
	GamePopucom:   "泡姆泡姆",
	GameNastPort:  "纳斯特港",
	GameCoreBlaze: "开拓芯",
}
