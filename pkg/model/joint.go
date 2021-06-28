package model

type ScheduleCourse struct {
	ID               string `json:"_key,omitempty"`   // increment unique key
	CourseID         string `json:"courseId"`         // document reference key
	ProviderID       string `json:"providerId"`       // document reference key
	CourseStartAt    int64  `json:"courseStartAt"`    // 課程開始時間
	CourseEndAt      int64  `json:"courseEndAt"`      // 課程結束時間
	MinConsumerLimit int    `json:"minConsumerLimit"` // 最小開課人數下限
	MaxConsumerLimit int    `json:"maxConsumerLimit"` // 最大開課人數上限
	Count            int    `json:"count"`            // 目前參加人數
	MeetingUrl       string `json:"meetingUrl"`       // 視訊平台連結
	Title            string `json:"title"`            // 課程標題
	Content          string `json:"content"`          // 課程內容
}
