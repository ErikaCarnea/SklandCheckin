package models

type CheckinRequest struct {
	GameID string `json:"gameId"`
}

type CheckinResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func (r CheckinResponse) GetCode() int {
	return r.Code
}
func (r CheckinResponse) GetMessage() string {
	return r.Message
}

type IsCheckinResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Checked int `json:"checked"`
			GameId  int `json:"gameId"`
		} `json:"list"`
	} `json:"data"`
	Timestamp string `json:"timestamp"`
}

func (r IsCheckinResponse) GetCode() int {
	return r.Code
}

func (r IsCheckinResponse) GetMessage() string {
	return r.Message
}
