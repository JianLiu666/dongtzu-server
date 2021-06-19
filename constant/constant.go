package constant

const (
	// Module

	Module_Initialization_Notyet  = iota // 尚未進行初始化
	Module_Initialization_Failed         // 初始化失敗
	Module_Initialization_Success        // 初始化成功
	Module_Initialization_Already        // 已經初始化過

	// ArangoDB

	ArangoDB_Success           // 操作成功
	ArangoDB_Invalid_Operation // 不允許的操作
	ArangoDB_Driver_Failed     // ArangoDB Driver 操作失敗

	// ZoomSDK

	ZoomSDK_Success       // 操作成功
	ZoomSDK_Driver_Failed // ZoomSDK Driver 操作失敗

	// LineSDK

	LineSDK_Success            // 操作成功
	LineSDK_ChannelID_NotFound // 找不到對應 ChannelID 的 Provider
	LineSDK_Request_Invalid    // Request 不合法
	LineSDK_Event_ParseFaild   // 解析 linebot.Event 失敗

	// Model.Appointment.Status

	Appointment_Status_Cancelled          // 取消預約(軟刪除)
	Appointment_Status_Exception          // 異常
	Appointment_Status_Unsend_MeetingUrl  // 尚未發送會議連結(表示預約成功, 課程尚未開始)
	Appointment_Status_Unsend_FeedbackUrl // 尚未發送回饋連結(表示課程尚未結束)
	Appointment_Status_Unverified         // 已結束(課程結束且已發連結, 尚未核銷)
	Appointment_Status_Varified           // 已核銷

	// Model.Consumer.LineStatus

	Consumer_LineStatus_Following   // 追蹤中
	Consumer_LineStatus_Unfollowing // 取消追蹤

	// Model.Provider.Status

	Provider_Status_TmpSave   // 暫存
	Provider_Status_Saved     // 確認送出
	Provider_Status_Auditing  // 審核中(Notion 代辦事項)
	Provider_Status_Audited   // 審核完成
	Provider_Status_Forbidden // 審核不通過, 黑名單
	Provider_Status_Unknown   // 例外處理
)
