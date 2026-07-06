package player

type StageInfo struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ZoneId      string `json:"zoneId"`
	DiffGroup   string `json:"diffGroup"`
	StageType   string `json:"stageType"`
	DangerLevel string `json:"dangerLevel"`
	ApCost      int    `json:"apCost"`
	Difficulty  string `json:"difficulty"`
}
