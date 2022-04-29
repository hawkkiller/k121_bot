package schedule

type Day struct {
	Caption string `json:"caption"`
	ID      int    `json:"id"`
	Pairs   []Pair `json:"pairs"`
}

type Schedule struct {
	ChatId int   `json:"-"`
	ID     int   `json:"id"`
	Days   []Day `json:"days"`
}

type Pair struct {
	Title          string `json:"title"`
	Auditory       string `json:"auditory"`
	ID             int    `json:"id"`
	AdditionalInfo string `json:"additional_info"`
}
