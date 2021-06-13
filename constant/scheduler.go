package constant

const (
	ApptStatus_Exception             = -1 // 異常
	ApptStatus_Unstarted_Unconfirmed = 0  // 尚未開始,且尚未被Scheduler確認
	ApptStatus_Unstarted_Confirmed   = 1  // 尚未開始,已被Scheduler確認
	ApptStatus_Starting              = 2  // 進行中
	ApptStatus_End_Unconfirmed       = 3  // 已結束,且尚未提供回饋
	ApptStatus_End_Confirmed         = 4  // 已結束,已提供回饋
	ApptStatus_Varified              = 5  // 已核銷
)
