package models

import "github.com/iancoleman/orderedmap"

type AttendanceInfo struct {
	Code      int            `json:"code"`
	Message   string         `json:"message"`
	Timestamp string         `json:"timestamp"`
	Data      AttendanceData `json:"data"`
}

// GetCode implements APIResponse.
func (a *AttendanceInfo) GetCode() int {
	return a.Code
}

// GetMessage implements APIResponse.
func (a *AttendanceInfo) GetMessage() string {
	return a.Message
}

type AttendanceData struct {
	CurrentTs       string                 `json:"currentTs"`
	Calendar        []calendar             `json:"calendar"`
	Records         []record               `json:"records"`
	ResourceInfoMap *orderedmap.OrderedMap `json:"resourceInfoMap"`
}

type calendar struct {
	ResourceId string `json:"resourceId"`
	Type       string `json:"type"`
	Count      int    `json:"count"`
	Available  bool   `json:"available"`
	Done       bool   `json:"done"`
}

type record struct {
	Ts         string `json:"ts"`
	ResourceId string `json:"resourceId"`
	Type       string `json:"type"`
	Count      int    `json:"count"`
}
