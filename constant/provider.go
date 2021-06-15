package constant

// Provider建置資料的狀態
const (
	ProviderStatusTmpSave   = 0 // 暫存
	ProviderStatusSaved     = 1 // 確認送出
	ProviderStatusAuditing  = 2 // 審核中(Notion 代辦事項)
	ProviderStatusAudited   = 3 // 審核完成
	ProviderStatusForbidden = 4 // 審核不通過, 黑名單
	ProviderStatusUnknown   = 5 // 例外處理
)
