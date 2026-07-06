package player

type PlayerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    data   `json:"data"`
}

func (p PlayerResponse) GetCode() int {
	return p.Code
}

func (p PlayerResponse) GetMessage() string {
	return p.Message
}

type data struct {
	CurrentTs  int64 `json:"currentTs"` // 获取数据时的时间戳
	ShowConfig struct {
		CharSwitch      bool `json:"charSwitch"`
		SkinSwitch      bool `json:"skinSwitch"`
		StandingsSwitch bool `json:"standingsSwitch"`
	} `json:"showConfig"`
	Status      status        `json:"status"`
	Medal       medal         `json:"medal"`       // 玩家奖章
	AssistChars []assistChars `json:"assistChars"` // 玩家助战
	Chars       []chars       `json:"chars"`
	Skins       []struct {
		Id string `json:"id"`
		Ts int64  `json:"ts"`
	} `json:"skins"`
	Building building `json:"building"` // 基建
	Recruit  [4]struct {
		StartTs  int64 `json:"startTs"`
		FinishTs int64 `json:"finishTs"`
		State    int   `json:"state"`
	} `json:"recruit"` // 公招
	Campaign                  campaign                          `json:"campaign"` // 剿灭
	Tower                     tower                             `json:"tower"`    // 保全派驻
	Rogue                     rogue                             `json:"rogue"`    // 肉鸽
	Routine                   routine                           `json:"routine"`  // 日常周常
	Activity                  []activity                        `json:"activity"`
	CharInfoMap               map[string]CharInfo               `json:"charInfoMap"`
	SkinInfoMap               map[string]SkinInfo               `json:"skinInfoMap"`
	StageInfoMap              map[string]StageInfo              `json:"stageInfoMap"`
	ActivityInfoMap           map[string]ActivityInfo           `json:"activityInfoMap"`
	TowerInfoMap              map[string]TowerInfo              `json:"towerInfoMap"`
	RogueInfoMap              map[string]RogueInfo              `json:"rogueInfoMap"`
	CampaignInfoMap           map[string]CampaignInfo           `json:"campaignInfoMap"`
	CampaignZoneInfoMap       map[string]CampaignZoneInfo       `json:"campaignZoneInfoMap"`
	EquipmentInfoMap          map[string]EquipmentInfo          `json:"equipmentInfoMap"`
	ManufactureFormulaInfoMap map[string]ManufactureFormulaInfo `json:"manufactureFormulaInfoMap"`
	CharAssets                []string                          `json:"charAssets"`
	SkinAssets                []string                          `json:"skinAssets"`
	CharAssetList             map[string][]string               `json:"charAssetList"`
	SkinAssetList             map[string][]string               `json:"skinAssetList"`
	ActivityBannerList        map[string][]string               `json:"activityBannerList"`
	BossRush                  []bossRush                        `json:"bossRush"`
	BannerList                []bannerList                      `json:"bannerList"`
	Sandbox                   []sandbox                         `json:"sandbox"`
}

type status struct {
	Uid    string `json:"uid"`   // 玩家UID
	Name   string `json:"name"`  // 玩家昵称
	Level  int    `json:"level"` // 玩家等级
	Avatar struct {
		Type string `json:"type"`
		Id   string `json:"id"`
		Url  string `json:"url"`
	} `json:"avatar"` // 玩家头像
	RegisterTs        int64  `json:"registerTs"`        // 玩家注册时间戳
	MainStageProgress string `json:"mainStageProgress"` // 玩家主线进程
	Secretary         struct {
		CharId string `json:"charId"`
		SkinId string `json:"skinId"`
	} `json:"secretary"` // 玩家助理干员
	Resume          string `json:"resume"`
	SubscriptionEnd int64  `json:"subscriptionEnd"`
	Ap              struct {
		Current              int   `json:"current"`              // 当前理智
		Max                  int   `json:"max"`                  // 最大理智
		LastApAddTime        int64 `json:"lastApAddTime"`        // 上次理智恢复时间戳
		CompleteRecoveryTime int64 `json:"completeRecoveryTime"` // 理智完全恢复时间戳
	} `json:"ap"` // 理智
	StoreTs      int64 `json:"storeTs"`
	LastOnlineTs int64 `json:"lastOnlineTs"` // 上次在线时间
	CharCnt      int   `json:"charCnt"`
	FurnitureCnt int   `json:"furnitureCnt"`
	SkinCnt      int   `json:"skinCnt"`
	Exp          struct {
		Current int `json:"current"`
		Max     int `json:"max"`
	} `json:"exp"` // 玩家经验值
}

type medal struct {
	Type              string `json:"type"`
	Template          string `json:"template"`
	TemplateMedalList []any  `json:"templateMedalList"`
	CustomMedalLayout []struct {
		Id  string `json:"id"`
		Pos []int  `json:"pos"`
	} `json:"customMedalLayout"`
	Total int `json:"total"`
}

type assistChars struct {
	CharId          string `json:"charId"`
	SkinId          string `json:"skinId"`
	Level           int    `json:"level"`
	EvolvePhase     int    `json:"evolvePhase"`
	PotentialRank   int    `json:"potentialRank"`
	SkillId         string `json:"skillId"`
	MainSkillLvl    int    `json:"mainSkillLvl"`
	SpecializeLevel int    `json:"specializeLevel"`
	Equip           struct {
		Id     string `json:"id"`
		Level  int    `json:"level"`
		Locked bool   `json:"locked"`
	} `json:"equip"`
}

type chars struct {
	CharId        string `json:"charId"`
	SkinId        string `json:"skinId"`
	Level         int    `json:"level"`
	EvolvePhase   int    `json:"evolvePhase"`
	PotentialRank int    `json:"potentialRank"`
	MainSkillLvl  int    `json:"mainSkillLvl"`
	Skills        []struct {
		Id              string `json:"id"`
		SpecializeLevel int    `json:"specializeLevel"`
	} `json:"skills"`
	Equip []struct {
		Id     string `json:"id"`
		Level  int    `json:"level"`
		Locked bool   `json:"locked"`
	} `json:"equip"`
	FavorPercent   int    `json:"favorPercent"`
	DefaultSkillId string `json:"defaultSkillId"`
	GainTime       int64  `json:"gainTime"`
	DefaultEquipId string `json:"defaultEquipId"`
	SortId         int    `json:"sortId"`
	Exp            int    `json:"exp"`
	Gold           int    `json:"gold"`
	Rarity         int    `json:"rarity"`
}

type building struct {
	TiredChars   []struct{}    `json:"tiredChars"`
	Powers       []power       `json:"powers"`
	Manufactures []manufacture `json:"manufactures"`
	Tradings     []trading     `json:"tradings"`
	Dormitories  []dormitory   `json:"dormitories"`
	Meeting      meeting       `json:"meeting"`
	Hire         hire          `json:"hire"`
	Training     training      `json:"training"`
	Labor        labor         `json:"labor"`
	Furniture    struct {
		Total int `json:"total"`
	} `json:"furniture"`
	Elevators []common `json:"elevators"`
	Corridors []common `json:"corridors"`
	Control   control  `json:"control"`
}

type common struct {
	SlotId    string `json:"slotId"`
	SlotState int    `json:"slotState"`
	Level     int    `json:"level"`
}

type power struct {
	SlotId string       `json:"slotId"`
	Level  int          `json:"level"`
	Chars  []powerChars `json:"chars"`
}
type powerChars struct {
	CharId        string `json:"charId"`
	Ap            int    `json:"ap"`
	LastApAddTime int64  `json:"lastApAddTime"`
	Index         int    `json:"index"`
	Bubble        struct {
		Normal struct {
			Add int   `json:"add"`
			Ts  int64 `json:"ts"`
		} `json:"normal"`
		Assist struct {
			Add int   `json:"add"`
			Ts  int64 `json:"ts"`
		} `json:"assist"`
	} `json:"bubble"`
	WorkTime int `json:"workTime"`
}

type manufacture struct {
	power
	CompleteWorkTime int64   `json:"completeWorkTime"`
	LastUpdateTime   int64   `json:"lastUpdateTime"`
	FormulaId        string  `json:"formulaId"`
	Capacity         int     `json:"capacity"`
	Weight           int     `json:"weight"`
	Complete         int     `json:"complete"`
	Remain           int     `json:"remain"`
	Speed            float64 `json:"speed"`
}

type trading struct {
	power
	CompleteWorkTime int64         `json:"completeWorkTime"`
	LastUpdateTime   int64         `json:"lastUpdateTime"`
	Strategy         string        `json:"strategy"`
	Stock            []interface{} `json:"stock"`
	StockLimit       int           `json:"stockLimit"`
}

type dormitory struct {
	power
	Comfort int `json:"comfort"`
}

type meeting struct {
	power
	Clue struct {
		Own               int       `json:"own"`
		Received          int       `json:"received"`
		DailyReward       bool      `json:"dailyReward"`
		NeedReceive       int       `json:"needReceive"`
		Board             [7]string `json:"board"`
		Sharing           bool      `json:"sharing"`
		ShareCompleteTime int64     `json:"shareCompleteTime"`
	} `json:"clue"`
	LastUpdateTime   int64 `json:"lastUpdateTime"`
	CompleteWorkTime int64 `json:"completeWorkTime"`
}

type hire struct {
	power
	State            int   `json:"state"`
	RefreshCount     int   `json:"refreshCount"`
	CompleteWorkTime int64 `json:"completeWorkTime"`
	SlotState        int   `json:"slotState"`
}

type training struct {
	SlotId  string `json:"slotId"`
	Level   int    `json:"level"`
	Trainee struct {
		CharId        string `json:"charId"`
		TargetSkill   int    `json:"targetSkill"`
		Ap            int    `json:"ap"`
		LastApAddTime int64  `json:"lastApAddTime"`
	} `json:"trainee"`
	Trainer struct {
		CharId        string `json:"charId"`
		Ap            int    `json:"ap"`
		LastApAddTime int64  `json:"lastApAddTime"`
	} `json:"trainer"`
	RemainPoint    float64 `json:"remainPoint"`
	Speed          float64 `json:"speed"`
	LastUpdateTime int64   `json:"lastUpdateTime"`
	RemainSecs     int     `json:"remainSecs"`
	SlotState      int     `json:"slotState"`
}

type labor struct {
	MaxValue       int   `json:"maxValue"`
	Value          int   `json:"value"`
	LastUpdateTime int64 `json:"lastUpdateTime"`
	RemainSecs     int   `json:"remainSecs"`
}

type control struct {
	power
	SlotState int `json:"slotState"`
}

type campaign struct {
	Records []struct {
		CampaignId string `json:"CampaignId"`
		MaxKills   int    `json:"MaxKills"`
	} `json:"records"`
	Reward struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"reward"`
}

type tower struct {
	Records []struct {
		TowerId string `json:"TowerId"`
		Best    int    `json:"best"`
	} `json:"records"`
	Reward struct {
		HigherItem struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"higherItem"`
		LowerItem struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"lowerItem"`
		TermTs int64 `json:"termTs"`
	} `json:"reward"`
}

type rogue struct {
	Records []rogueRecord `json:"records"`
}

type rogueRecord struct {
	RogueId  string `json:"rogueId"`
	RelicCnt int    `json:"relicCnt"`
	Bank     struct {
		Current int `json:"current"`
		Record  int `json:"record"`
	} `json:"bank"`
	ClearTime int `json:"clearTime"`
	BpLevel   int `json:"bpLevel"`
	Medal     struct {
		Total   int `json:"total"`
		Current int `json:"current"`
	} `json:"medal"`
}

type routine struct {
	Daily struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"daily"`
	Weekly struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"weekly"`
}

type activity struct {
	ActId        string         `json:"actId"`
	ActReplicaId string         `json:"actReplicaId"`
	Zones        []activityZone `json:"zones"`
}

type activityZone struct {
	ZoneId        string `json:"zoneId"`
	ZoneReplicaId string `json:"zoneReplicaId"`
	ClearedStage  int    `json:"clearedStage"`
	TotalStage    int    `json:"totalStage"`
}

type bossRush struct {
	Id     string `json:"id"`
	Record struct {
		Played     bool   `json:"played"`
		StageId    string `json:"stageId"`
		Difficulty string `json:"difficulty"`
	} `json:"record"`
	PicUrl string `json:"picUrl"`
}

type bannerList struct {
	Id        string `json:"id"`
	SortId    int    `json:"sortId"`
	ImgUrl    string `json:"imgUrl"`
	Link      string `json:"link"`
	StartAtTs string `json:"startAtTs"`
	EndAtTs   string `json:"endAtTs"`
	Status    int    `json:"status"`
}

type sandbox struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	MaxDay          int    `json:"maxDay"`
	MaxDayChallenge int    `json:"maxDayChallenge"`
	MainQuest       int    `json:"mainQuest"`
	SubQuest        []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Done bool   `json:"done"`
	} `json:"subQuest"`
	BaseLv     int    `json:"baseLv"`
	UnlockNode int    `json:"unlockNode"`
	EnemyKill  int    `json:"enemyKill"`
	CreateRift int    `json:"createRift"`
	FixRift    []int  `json:"fixRift"`
	PicUrl     string `json:"picUrl"`
}
