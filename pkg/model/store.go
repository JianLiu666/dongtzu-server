package model

type Appointment struct {
	ID             string `json:"_key"`           // increment unique key
	StartTimestamp int64  `json:"startTimestamp"` // 預約開始時間
	EndTimestamp   int64  `json:"endTimestamp"`   // 預約結束時間
	ProviderID     string `json:"providerId"`     // 預約提供者識別碼
	ProviderLineID string `json:"providerLineId"` // 預約提供者的 Line 官方帳號識別碼
	ConsumerID     string `json:"consumerId"`     // 消費者識別碼
	ConsumerLineID string `json:"consumerLineId"` // 消費者的 Line 官方帳號識別碼
	FeedbackID     string `json:"feedbackId"`     // 意見回饋識別碼
	Note           string `json:"note"`           // 備註 (json format)
	MeetingUrl     string `json:"meetingUrl"`     // 會議連結
	Status         int64  `json:"status"`         // -1:異常, 0:尚未開始,未發連結, 1:尚未開始,已發連結, 2:進行中, 3:已結束,未提供回饋, 4:已結束,已提供回饋, 5:已核銷
}
