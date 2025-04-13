package models

type AccountInfo struct {
	Phone    string `yaml:"phone"`
	Password string `yaml:"password"`
}

type AppConfig struct {
	Account struct {
		Phone    string
		Password string
	}
}

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
	Code int `json:"code"`
	Data struct {
		List []struct {
			AppCode     string    `json:"appCode"`
			BindingList []Binding `json:"bindingList"`
		} `json:"list"`
	} `json:"data"`
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
