package constant

// Module

const (
	Module_Initialization_Notyet  = iota // 尚未進行初始化
	Module_Initialization_Failed         // 初始化失敗
	Module_Initialization_Success        // 初始化成功
	Module_Initialization_Already        // 已經初始化過
)
