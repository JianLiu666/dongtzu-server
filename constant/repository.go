package constant

// ArangoDB

const (
	ArangoDB_Success           = iota // 操作成功
	ArangoDB_Invalid_Operation        // 不允許的操作
	ArangoDB_Driver_Failed            // ArangoDB Driver 操作失敗
)

// ZoomSDK

const (
	ZoomSDK_Success       = iota // 操作成功
	ZoomSDK_Driver_Failed        // ZoomSDK Driver 操作失敗
)

// LineSDK

const (
	LineSDK_Success            = iota // 操作成功
	LineSDK_ChannelID_NotFound        // 找不到對應 ChannelID 的 Provider
	LineSDK_Request_Invalid           // Request 不合法
	LineSDK_Event_ParseFaild          // 解析 linebot.Event 失敗
)
