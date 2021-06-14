package constant

const (
	ApptStatus_Exception = -2 // 異常
	ApptStatus_Cancelled = -1 // 取消預約(軟刪除)
	ApptStatus_Unstarted = 0  // 預約成功,尚未開始(預設)
	ApptStatus_Starting  = 1  // 進行中
	ApptStatus_End       = 2  // 已結束
	ApptStatus_Varified  = 3  // 已核銷
)
