package models

type ExaResult struct {
	Code      int     `json:"code"`
	Message   string  `json:"message"`
	Timestamp string  `json:"timestamp"`
	Data      exaData `json:"data"`
}

// GetCode implements APIResponse.
func (e ExaResult) GetCode() int {
	return e.Code
}

// GetMessage implements APIResponse.
func (e ExaResult) GetMessage() string {
	return e.Message
}

type exaData struct {
	Achievements []exaItem  `json:"achievements"`
	Awards       []exaAward `json:"awards"`
	Privacy      struct {
		CardOn   bool `json:"cardOn"`
		DetailOn bool `json:"detailOn"`
	}
	Stats exaStats `json:"stats"`
}

type exaStats struct {
	Hidden struct {
		Acquired int `json:"acquired"`
		Total    int `json:"total"`
	} `json:"hidden"`
	Highest string `json:"highest"`
	Routine struct {
		Acquired int `json:"acquired"`
		Total    int `json:"total"`
	}
	Uid string `json:"uid"`
}

type exaAward struct {
	Avatar   string `json:"avatar"`
	Claimed  bool   `json:"claimed"`
	Id       string `json:"id"`
	Unlocked bool   `json:"unlocked"`
}

type exaItem struct {
	Acquired bool    `json:"acquired"`
	Desc     string  `json:"desc"`
	Hidden   bool    `json:"hidden"`
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Rate     float32 `json:"rate"`
	Ts       string  `json:"ts"`
}
