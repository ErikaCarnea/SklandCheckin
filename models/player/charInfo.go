package player

type CharInfo struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	NationId          string `json:"nationId"`
	GroupId           string `json:"groupId"`
	DisplayNumber     string `json:"displayNumber"`
	Rarity            int    `json:"rarity"`
	Profession        string `json:"profession"`
	SubProfessionId   string `json:"subProfessionId"`
	SubProfessionName string `json:"subProfessionName"`
	Appellation       string `json:"appellation"`
	SortId            int    `json:"sortId"`
}
