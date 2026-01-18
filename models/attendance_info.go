package models

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
	CurrentTs       string               `json:"currentTs"`
	Calendar        []calendar           `json:"calendar"`
	Records         []record             `json:"records"`
	ResourceInfoMap map[int]ResourceInfo `json:"resourceInfoMap"`
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

type ResourceInfo struct {
	BuildingProductList []BuildingProduct `json:"buildingProductList"`
	ClassifyType        string            `json:"classifyType"`
	Id                  string            `json:"id"`
	Name                string            `json:"name"`
	OtherSource         []string          `json:"otherSource"`
	Rarity              int               `json:"rarity"`
	SortId              int               `json:"sortId"`
	StageDropList       []StageDrop       `json:"stageDropList"`
	Type                string            `json:"type"`
}

type BuildingProduct struct {
	FormulaId string `json:"formulaId"`
	RoomType  string `json:"roomType"`
}

type StageDrop struct {
	OccPer   float32 `json:"occPer"`
	StatgeId string  `json:"stageId"`
}
