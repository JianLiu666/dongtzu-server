package constant

const (
	ApptStatus_Exception          = -2 // 異常
	ApptStatus_Cancelled          = -1 // 取消預約(軟刪除)
	ApptStatus_Unsend_MeetingUrl  = 0  // 尚未發送會議連結(表示預約成功, 課程尚未開始)
	ApptStatus_Unsend_FeedbackUrl = 1  // 尚未發送回饋連結(表示課程尚未結束)
	ApptStatus_Unverified         = 2  // 已結束(課程結束且已發連結, 尚未核銷)
	ApptStatus_Varified           = 3  // 已核銷
)
