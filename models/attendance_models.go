package models

type AttendanceRequest struct {
	Uid    string `json:"uid"`
	GameId string `json:"channelMasterId"`
}

type AttendanceResult struct {
	Code int `json:"code"`
	Data struct {
		Timestamp string `json:"ts"`
		Awards    []struct {
			Resource struct {
				Name string `json:"name"`
			} `json:"resource"`
			Count int    `json:"count"`
			Type  string `json:"type"`
		} `json:"awards"`
	} `json:"data"`
	Message string `json:"message"`
}

func (r AttendanceResult) GetCode() int {
	return r.Code
}
func (r AttendanceResult) GetMessage() string {
	return r.Message
}

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
