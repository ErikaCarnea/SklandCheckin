package models

type AttendanceRequest struct {
	Uid    string `json:"uid"`
	GameId string `json:"channelMasterId"`
}

type AttendanceResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Timestamp       string           `json:"ts"`
		Awards          []Award          `json:"awards"`
		ResourceInfoMap map[int]Resource `json:"resourceInfoMap"`
		TomorrowAwards  []TomorrowAward  `json:"tomorrowAwards"`
	} `json:"data"`
	Timestamp string `json:"timestamp"`
}

func (r AttendanceResult) GetCode() int {
	return r.Code
}
func (r AttendanceResult) GetMessage() string {
	return r.Message
}

type Award struct {
	Resource struct {
		Name string `json:"name"`
	} `json:"resource"`
	Count int    `json:"count"`
	Type  string `json:"type"`
}

type TomorrowAward struct {
	Count    int      `json:"count"`
	Resource Resource `json:"resource"`
	Type     string   `json:"type"`
}

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
	CurrentTs       string           `json:"currentTs"`
	Calendar        []calendar       `json:"calendar"`
	Records         []record         `json:"records"`
	ResourceInfoMap map[int]Resource `json:"resourceInfoMap"`
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

type Resource struct {
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
