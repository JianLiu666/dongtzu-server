package constant

// Model.Appointment.Status

const (
	Appointment_Status_Cancelled          = iota // 取消預約(軟刪除)
	Appointment_Status_Exception                 // 異常
	Appointment_Status_Unsend_MeetingUrl         // 尚未發送會議連結(表示預約成功, 課程尚未開始)
	Appointment_Status_Unsend_FeedbackUrl        // 尚未發送回饋連結(表示課程尚未結束)
	Appointment_Status_Unverified                // 已結束(課程結束且已發連結, 尚未核銷)
	Appointment_Status_Varified                  // 已核銷
)

// Model.Consumer.LineStatus

const (
	Consumer_LineStatus_Following   = iota // 追蹤中
	Consumer_LineStatus_Unfollowing        // 取消追蹤
)

// Model.Provider.Status

const (
	Provider_Status_TmpSave   = iota // 暫存
	Provider_Status_Saved            // 確認送出
	Provider_Status_Auditing         // 審核中(Notion 代辦事項)
	Provider_Status_Audited          // 審核完成
	Provider_Status_Forbidden        // 審核不通過, 黑名單
	Provider_Status_Unknown          // 例外處理
)
