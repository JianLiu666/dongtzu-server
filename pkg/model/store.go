package model

type ClassRoom struct {
	ID        string `json:"_key"`       // increment unique key
	SchduleID string `json:"scheduleId"` // document reference key
	Title     string `json:"title"`      // 課程標題
	Content   string `json:"content"`    // 課程內容
}

type Schedule struct {
	ID               string `json:"_key"`             // increment unique key
	ProviderID       string `json:"providerId"`       // document reference key
	StartTimestamp   int64  `json:"startTimestamp"`   // 預約開始時間
	EndTimestamp     int64  `json:"endTimestamp"`     // 預約結束時間
	MinConsumerLimit int64  `json:"minConsumerLimit"` // 最小開課人數下限
	MaxConsumerLimit int64  `json:"maxConsumerLimit"` // 最大開課人數上限
	Count            int64  `json:"count"`            // 目前參加人數
	MeetingUrl       string `json:"meetingUrl"`       // 視訊平台連結
}

type Appointment struct {
	ID             string `json:"_key"`           // increment unique key
	ProviderID     string `json:"providerId"`     // document reference key
	ScheduleID     string `json:"scheduleId"`     // document reference key
	ConsumerID     string `json:"consumerId"`     // document reference key
	FeedbackID     string `json:"feedbackId"`     // document reference key
	StartTimestamp int64  `json:"startTimestamp"` // 預約開始時間
	EndTimestamp   int64  `json:"endTimestamp"`   // 預約結束時間
	Note           string `json:"note"`           // 備註 (json format)
	Status         int64  `json:"status"`         // -1:異常, 0:尚未開始,未發連結, 1:尚未開始,已發連結, 2:進行中, 3:已結束,未提供回饋, 4:已結束,已提供回饋, 5:已核銷
}

type Feedback struct {
	ID            string `json:"_key"`          // increment unique key
	SchduleID     string `json:"scheduleId"`    // document reference key
	AppointmentID string `json:"appointmentId"` // document reference key
	Title         string `json:"title"`         // 回饋標題
	Content       string `json:"content"`       // 回饋內容
}
