package models

import "fmt"

type Binding struct {
	Uid             string `json:"uid"`
	ChannelMasterId string `json:"channelMasterId"`
	ChannelName     string `json:"channelName"`
	NickName        string `json:"nickName"`
}

func (b Binding) ToString() string {
	return fmt.Sprintf("[%s] UID:%s %s", b.ChannelName, b.Uid, b.NickName)
}

type BindingResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			AppCode     string    `json:"appCode"`
			AppName     string    `json:"appName"`
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
