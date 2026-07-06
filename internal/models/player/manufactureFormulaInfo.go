package player

type ManufactureFormulaInfo struct {
	Id        string `json:"id"`
	ItemId    string `json:"itemId"`
	Count     int    `json:"count"`
	Weight    int    `json:"weight"`
	Costs     []any  `json:"costs"`
	CostPoint int    `json:"costPoint"`
}
