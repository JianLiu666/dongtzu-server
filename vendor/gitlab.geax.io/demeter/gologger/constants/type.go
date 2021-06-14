package constants

// ServiceCode 錯誤代碼型別，呼叫 gologger 方法若傳入此型別，將自動填入服務代碼欄位
type ServiceCode int

// ServiceCode list
const (
	Roomservice     ServiceCode = 100
	Gamemaster      ServiceCode = 101
	Gamecontroller  ServiceCode = 102
	GameHallBackend ServiceCode = 103
	NodePublic      ServiceCode = 104
	Scheduler       ServiceCode = 105
)
