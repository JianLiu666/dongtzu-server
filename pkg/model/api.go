package model

type ErrRes struct {
	Code       string `json:"code"`    // 系統自定義錯誤代碼
	StatusCode string `json:"status"`  // http 狀態碼
	Message    string `json:"message"` // 錯誤訊息
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
	LineAtID    string `json:"LineAtID"`
	CountryCode string `json:"countryCode"`
	PhoneNum    string `json:"phoneNum"`
	GmailAddr   string `json:"gmailAddr"`
	InivteCode  string `json:"inivteCode,omitempty"`
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
	Status      int    `json:"status"`
}

type RegisterOrUpdateProviderRes struct {
	StatusCode string `json:"status"`
}
