package model

type ErrRes struct {
	Code       string `json:"code"`       // 系統自定義錯誤代碼
	StatusCode string `json:"statusCode"` // http 狀態碼
	Message    string `json:"message"`    // 錯誤訊息
}

type CreateUpdateDeleteRes struct {
	StatusCode string `json:"statusCode"`
}

type Pagination struct {
	TotalCounts int `json:"totalCounts"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
}

type GetProviderInfoRes struct {
	StatusCode string    `json:"statusCode"`
	Data       *Provider `json:"data"`
}

type RegisterProviderReq struct {
	LineUserID  string `json:"lineUserId"`
	RealName    string `json:"realName"`
	LineID      string `json:"lineId"`
	LineAtName  string `json:"lineAtName"`
	LineAtID    string `json:"lineAtID"`
	CountryCode string `json:"countryCode"`
	PhoneNum    string `json:"phoneNum"`
	GmailAddr   string `json:"gmailAddr"`
	InviteCode  string `json:"inviteCode,omitempty"`
	MemeberTerm bool   `json:"memeberTerm,omitempty"`
	PrivacyTerm bool   `json:"privacyTerm,omitempty"`
	Status      int    `json:"status,omitempty"`
}

type UpdateProviderInfoReq struct {
	RealName    string `json:"realName"`
	LineID      string `json:"lineId"`
	LineAtName  string `json:"lineAtName"`
	LineAtID    string `json:"LineAtID"`
	CountryCode string `json:"countryCode"`
	PhoneNum    string `json:"phoneNum"`
	GmailAddr   string `json:"gmailAddr"`
	GUUID       string `json:"guuid"`
	GToken      string `json:"gToken"`
	GRawData    string `json:"gRawData"`
	Status      int    `json:"status"`
}

type RegisterOrUpdateProviderRes struct {
	StatusCode string `json:"statusCode"`
}

type IncomeSummary struct {
	SumOfPayments      int `json:"sumOfPayments"`
	LastMonthIncome    int `json:"lastMonthIncome"`
	CurrentMonthIncome int `json:"currentMonthIncome"`
}

type ProviderIncomeSummaryRes struct {
	StatusCode string        `json:"statusCode"`
	Data       IncomeSummary `json:"data"`
}

type EventSchedule struct {
	Start        int64  `json:"start"`
	End          int64  `json:"end"`
	ScheduleType int    `json:"scheduleType"` // 0 schedule / 1 appointment
	Title        string `json:"title"`
	Content      string `json:"content"`
	Count        int    `json:"count"`
}

type ProviderEventScheduleRes struct {
	StatusCode string          `json:"statusCode"`
	Data       []EventSchedule `json:"data"`
}

type ProviderMonthReceiptList struct {
	StatusCode string         `json:"statusCode"`
	Data       []MonthReceipt `json:"data"`
	Meta       Pagination     `json:"meta"`
}

type SvcProduct struct {
	ID              string `json:"_key,omitempty"`  // increment unique key
	ProviderID      string `json:"providerId"`      // document reference key
	CountPerPack    int    `json:"countPerPack"`    // 一包多少堂
	Price           int    `json:"price"`           // 一堂多少價格
	ExpiredDuration int64  `json:"expiredDuration"` // 多久過期
}

type CreateOrUpdateServiceProductsReq struct {
	ReqList []SvcProduct `json:"reqList"`
}

type ProviderServiceProductListRes struct {
	StatusCode string           `json:"statusCode"`
	Data       []ServiceProduct `json:"data"`
}

type ProviderScheduleListRes struct {
	StatusCode string     `json:"statusCode"`
	Data       []Schedule `json:"data"`
	Meta       Pagination `json:"meta"`
}

type CreateServiceScheduleReq struct {
	CourseStartAt    int64  `json:"courseStartAt"`    // 課程開始時間
	CourseEndAt      int64  `json:"courseEndAt"`      // 課程結束時間
	MinConsumerLimit int    `json:"minConsumerLimit"` // 最小開課人數下限
	MaxConsumerLimit int    `json:"maxConsumerLimit"` // 最大開課人數上限
	Title            string `json:"title"`            // 課程標題
	Content          string `json:"content"`          // 課程內容
}

type CreateServiceScheduleRuleReq struct {
	CourseStartAt       int64  `json:"courseStartAt"`       // 課程開始時間
	CourseEndAt         int64  `json:"courseEndAt"`         // 課程結束時間
	CycleStartAt        int64  `json:"cycleStartAt"`        // 週期開始時間
	CycleEndAt          int64  `json:"cycleEndAt"`          // 週期結束時間 -> 配合結束類型看
	CycleEndType        int    `json:"cycleEndType"`        // 結束類型：0幾次、-1時間、-2永遠不停
	CycleRepeatedAmount int    `json:"cycleRepeatedAmount"` // 搭配週期間隔數量、單位看
	CycleDiffAmount     int    `json:"cycleDiffAmount"`     // 週期間隔數量
	CycleDiffUnit       int    `json:"cycleDiffUnit"`       // 週期間隔單位(0-天、1-週、2-月、3-年)
	MinConsumerLimit    int    `json:"minConsumerLimit"`    // 最小開課人數下限
	MaxConsumerLimit    int    `json:"maxConsumerLimit"`    // 最大開課人數上限
	Title               string `json:"title"`               // 課程標題
	Content             string `json:"content"`             // 課程內容
}
