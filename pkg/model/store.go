package model

type ZoomAccount struct {
	UserID    string `json:"userId"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type Provider struct {
	ID                  string `json:"_key,omitempty"`      // increment unique key
	LineUserID          string `json:"lineUserId"`          // uniq Line userID (e.g. U123xxxxxxxxdef)
	RealName            string `json:"realName"`            // 中文姓名(真實)
	LineAtName          string `json:"lineAtName"`          // 申請Line官方帳號名稱
	LineAtID            string `json:"lineAtID"`            // 申請Line官方帳號ID(Todo 暫不實作)
	LineAtChannelID     string `json:"lineAtChannelId"`     // 申請Line官方帳號 ChannelID
	LineAtChannelSecret string `json:"lineAtChannelSecret"` // 申請Line官方帳號 ChannelSecret
	LineAtAccessToken   string `json:"lineAtAccessToken"`   // 申請Line官方帳號 AccessToken
	CountryCode         string `json:"countryCode"`         // 手機國碼(Todo 暫不實作)
	LineID              string `json:"lineID"`              // 聯絡的個人Line ID
	PhoneNum            string `json:"phoneNum"`            // 手機號碼
	ConfirmedPhoneNum   string `json:"confirmedPhoneNum"`   // 認證過的手機號碼(Todo 二階段手機驗證)
	GmailAddr           string `json:"gamilAddr"`           // Gamil
	ConfirmedGmailAddr  string `json:"confirmedGmailAddr"`  // 認證過的Gamil(Todo 發認證信)
	GCalSync            bool   `json:"gCalSync"`            // Google Calendar 授權成功
	InviteCode          string `json:"inviteCode"`          // 企業用戶，業務推廣碼
	MemeberTerm         bool   `json:"memeberTerm"`         // 會員條款
	PrivacyTerm         bool   `json:"privacyTerm"`         // 隱私全條款
	Status              int    `json:"status"`              // 狀態 0: 暫存, 1: 確認送出, 2: 審核中, 3: 審核完成, 4: 審核不通過, 5: 例外處理
	CreatedAt           int    `json:"createdAt"`           // 創建時間
	// Todo 應該還有個google calendar 授權成功拿到的token
}

type Consumer struct {
	ID                      string `json:"_key,omitempty"`          // increment unique key
	ProviderID              string `json:"providerId"`              // document reference key
	LineUserID              string `json:"lineUserId"`              // Line UserId
	LineFollowingStatus     int    `json:"lineFollowingStatus"`     // Line following status
	ProviderLineAtChannelID string `json:"providerLineAtChannelId"` // Provider Line官方帳號 ChannelID
}

type ClassRoom struct {
	ID        string `json:"_key,omitempty"` // increment unique key
	SchduleID string `json:"scheduleId"`     // document reference key
	Title     string `json:"title"`          // 課程標題
	Content   string `json:"content"`        // 課程內容
}

type Schedule struct {
	ID               string `json:"_key,omitempty"`   // increment unique key
	ProviderID       string `json:"providerId"`       // document reference key
	StartTimestamp   int64  `json:"startTimestamp"`   // 預約開始時間
	EndTimestamp     int64  `json:"endTimestamp"`     // 預約結束時間
	MinConsumerLimit int    `json:"minConsumerLimit"` // 最小開課人數下限
	MaxConsumerLimit int    `json:"maxConsumerLimit"` // 最大開課人數上限
	Count            int    `json:"count"`            // 目前參加人數
	MeetingUrl       string `json:"meetingUrl"`       // 視訊平台連結
}

type Appointment struct {
	ID             string `json:"_key,omitempty"` // increment unique key
	ProviderID     string `json:"providerId"`     // document reference key
	ScheduleID     string `json:"scheduleId"`     // document reference key
	ConsumerID     string `json:"consumerId"`     // document reference key
	FeedbackID     string `json:"feedbackId"`     // document reference key
	ConsumerLineID string `json:"consumerLineId"` // Consumer Line UserId
	StartTimestamp int64  `json:"startTimestamp"` // 預約開始時間
	EndTimestamp   int64  `json:"endTimestamp"`   // 預約結束時間
	Note           string `json:"note"`           // 備註 (json format)
	Status         int    `json:"status"`         // -1:異常, 0:尚未開始,未發連結, 1:尚未開始,已發連結, 2:進行中, 3:已結束,未提供回饋, 4:已結束,已提供回饋, 5:已核銷
}

type Feedback struct {
	ID            string `json:"_key,omitempty"` // increment unique key
	SchduleID     string `json:"scheduleId"`     // document reference key
	AppointmentID string `json:"appointmentId"`  // document reference key
	Title         string `json:"title"`          // 回饋標題
	Content       string `json:"content"`        // 回饋內容
}
