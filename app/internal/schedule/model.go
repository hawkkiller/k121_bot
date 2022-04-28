package schedule

type Day struct {
	Caption string `json:"caption"`
	Pairs   []Pair `json:"pairs"`
}

type Schedule struct {
	ChatId int64 `json:"-"`
	Days   []Day `json:"days"`
}

type Pair struct {
	Title          string `json:"title"`
	Auditory       string `json:"auditory"`
	AdditionalInfo string `json:"additional_info"`
}
