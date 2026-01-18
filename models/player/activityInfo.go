package player

type ActivityInfo struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	StartTime       int64    `json:"startTime"`
	EndTime         int64    `json:"endTime"`
	IsReplicate     bool     `json:"isReplicate"`
	Type            string   `json:"type"`
	DropItemIds     []string `json:"dropItemIds"`
	ShopGoodItemIds []string `json:"shopGoodItemIds"`
	FavorUpList     []string `json:"favorUpList"`
	PicUrl          string   `json:"picUrl"`
}
