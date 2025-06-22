package models

type PopucomResult struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
	Data      popucomData `json:"data"`
}

type popucomData struct {
	Achievements []popucomItem  `json:"achievements"`
	Awards       []popucomAward `json:"awards"`
	Privacy      struct {
		CardOn   bool `json:"cardOn"`
		DetailOn bool `json:"detailOn"`
	} `json:"privacy"`
	Stats popucomStats `json:"stats"`
	User  struct {
		Avatar      string `json:"avatar"`
		AvartarCode int    `json:"avartarCode"`
		Nickname    string `json:"nickname"`
	}
}

type popucomStats struct {
	Achievements struct {
		Acquired int `json:"acquired"`
		Total    int `json:"total"`
	} `json:"achievements"`
	HasGame  bool `json:"hasGame"`
	Plat     int  `json:"plat"`
	Playtime int  `json:"playtime"`
	Stickers struct {
		Acquired int `json:"acquired"`
		Total    int `json:"total"`
	}
	Uid     string `json:"uid"`
	Yuanbao struct {
		Acquired int `json:"acquired"`
		Total    int `json:"total"`
	}
}

type popucomAward struct {
	Avatar string `json:"avatar"`
	Id     int    `json:"id"`
	Level  int    `json:"level"`
	Status int    `json:"status"`
}

type popucomItem struct {
	Counter     any     `json:"counter"`
	Desc        string  `json:"desc"`
	HasAcquired bool    `json:"hasAcquired"`
	Hidden      bool    `json:"hidden"`
	Icon        string  `json:"icon"`
	Id          string  `json:"id"`
	Level       int     `json:"level"`
	Name        string  `json:"name"`
	Rate        float32 `json:"rate"`
	Ts          string  `json:"ts"`
}
