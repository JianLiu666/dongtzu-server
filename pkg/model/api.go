package model

type ErrRes struct {
	Code       string `json:"code"`    // 系統自定義錯誤代碼
	StatusCode string `json:"status"`  // http 狀態碼
	Message    string `json:"message"` // 錯誤訊息
}

type CreateUpdateDeleteRes struct {
	StatusCode string `json:"status"`
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
	StatusCode string `json:"status"`
}

type IncomeSummary struct {
	SumOfPayments      int `json:"sumOfPayments"`
	LastMonthIncome    int `json:"lastMonthIncome"`
	CurrentMonthIncome int `json:"currentMonthIncome"`
}

type ProviderIncomeSummaryRes struct {
	StatusCode string        `json:"status"`
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
	StatusCode string          `json:"status"`
	Data       []EventSchedule `json:"data"`
}

type ProviderMonthReceiptList struct {
	StatusCode string         `json:"status"`
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
	StatusCode string           `json:"status"`
	Data       []ServiceProduct `json:"data"`
}

type ProviderScheduleListRes struct {
	StatusCode string     `json:"status"`
	Data       []Schedule `json:"data"`
	Meta       Pagination `json:"meta"`
}
