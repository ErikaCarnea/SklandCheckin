package models

type LoginResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type GrantResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Data    struct {
		Code string `json:"code"`
	} `json:"data"`
}

func (g GrantResult) GetCode() int {
	return g.Status
}

func (g GrantResult) GetMessage() string {
	return g.Message
}

type CredResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Cred   string `json:"cred"`
		Token  string `json:"token"`
		UserId string `json:"userId"`
	} `json:"data"`
}

type Binding struct {
	Uid             string `json:"uid"`
	ChannelMasterId string `json:"channelMasterId"`
	ChannelName     string `json:"channelName"`
	NickName        string `json:"nickName"`
}

type BindingResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			AppCode     string    `json:"appCode"`
			BindingList []Binding `json:"bindingList"`
		} `json:"list"`
	} `json:"data"`
}

func (r BindingResult) GetCode() int {
	return r.Code
}

func (r BindingResult) GetMessage() string {
	return r.Message
}

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

// GetCode LoginResult 实现接口
func (r *LoginResult) GetCode() int {
	return r.Status
}

func (r *LoginResult) GetMessage() string {
	return r.Message
}

// GetCode CredResult 实现接口
func (r *CredResult) GetCode() int {
	return r.Code
}

func (r *CredResult) GetMessage() string {
	return r.Message
}

type SendCodeResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

func (r *SendCodeResult) GetCode() int {
	return r.Status
}

func (r *SendCodeResult) GetMessage() string {
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

func (r CheckinResponse) GetCode() int       { return r.Code }
func (r CheckinResponse) GetMessage() string { return r.Message }
