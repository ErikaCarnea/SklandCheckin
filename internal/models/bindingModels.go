package models

import (
	"encoding/json"
	"fmt"
)

type Binding struct {
	ChannelMasterId string      `json:"channelMasterId"`
	ChannelName     string      `json:"channelName"`
	DefaultRole     DefaultRole `json:"defaultRole"`
	GameId          int         `json:"gameId"`
	GameName        string      `json:"gameName"`
	IsDefault       bool        `json:"isDefault"`
	IsDelete        bool        `json:"isDelete"`
	IsOfficial      bool        `json:"isOfficial"`
	NickName        string      `json:"nickName"`
	Roles           []Role      `json:"roles"`
	Uid             string      `json:"uid"`
}

func (b Binding) ToString() string {
	return fmt.Sprintf("[%s] UID:%s %s", b.ChannelName, b.Uid, b.NickName)
}

type DefaultRole struct {
	ServerId   string `json:"serverId"`
	RoleId     string `json:"roleId"`
	NickName   string `json:"nickName"`
	Level      int    `json:"level"`
	IsDefault  bool   `json:"isDefault"`
	IsBanned   bool   `json:"isBanned"`
	ServerType string `json:"serverType"`
	ServerName string `json:"serverName"`
}

type BindingResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			AppCode     string    `json:"appCode"`
			AppName     string    `json:"appName"`
			BindingList []Binding `json:"bindingList"`
			DefaultUid  string    `json:"defaultUid"`
		} `json:"list"`
		ServerDefaultBinding json.RawMessage `json:"serverDefaultBinding"`
	} `json:"data"`
}

func (r BindingResult) GetCode() int {
	return r.Code
}
func (r BindingResult) GetMessage() string {
	return r.Message
}

type Role struct {
	ServerId   string `json:"serverId"`
	RoleId     string `json:"roleId"`
	NickName   string `json:"nickName"`
	Level      int    `json:"level"`
	IsDefault  bool   `json:"isDefault"`
	IsBanned   bool   `json:"isBanned"`
	ServerType string `json:"serverType"`
	ServerName string `json:"serverName"`
}

func (r Role) ToString() string {
	return fmt.Sprintf("[%s] UID:%s %s", r.ServerName, r.RoleId, r.NickName)
}
