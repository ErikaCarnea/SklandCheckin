package models

type EndfieldResult struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Timestamp string           `json:"timestamp"`
	Data      EndfieldSignData `json:"data"`
}

func (e EndfieldResult) GetCode() int {
	return e.Code
}
func (e EndfieldResult) GetMessage() string {
	return e.Message
}

type EndfieldSignData struct {
	Ts               string                      `json:"ts"`
	AwardIds         []AwardIds                  `json:"awardIds"`
	ResourceInfoMap  map[string]EndfieldResource `json:"resourceInfoMap"`
	TomorrowAwardIds []AwardIds                  `json:"tomorrowAwardIds"`
}

type AwardIds struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type EndfieldResource struct {
	Id    string `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
}
