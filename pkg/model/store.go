package model

type Appointment struct {
	ID               string `json:"_key,omitempty"`   // increment unique key
	ProviderID       string `json:"providerId"`       // document reference key
	ScheduleID       string `json:"scheduleId"`       // document reference key
	ConsumerID       string `json:"consumerId"`       // document reference key
	FeedbackID       string `json:"feedbackId"`       // document reference key
	PaymentMethodID  string `json:"paymentMethodId"`  // document reference key
	MonthReceiptID   string `json:"monthReceiptId"`   // document reference key
	ConsumerLineID   string `json:"consumerLineId"`   // Consumer Line UserId
	PaymentMethodFee int32  `json:"paymentMethodFee"` // 出款結算透過的金流所扣的錢
	Status           int    `json:"status"`           // -1:異常, 0:尚未開始,未發連結, 1:尚未開始,已發連結, 2:進行中, 3:已結束,未提供回饋, 4:已結束,已提供回饋, 5:已核銷
	Note             string `json:"note"`             // 備註 (json format)
	CourseStartAt    int64  `json:"courseStartAt"`    // 課程開始時間
	CourseEndAt      int64  `json:"courseEndAt"`      // 課程結束時間
	CreatedAt        int64  `json:"createdAt"`        //
	UpdatedAt        int64  `json:"updatedAt"`        //
	DeletedAt        int64  `json:"deletedAt"`        //
	RescheduledAt    int64  `json:"rescheduledAt"`    // Todo 改期，還漏一些相關要同步狀態的column設計
}

type Course struct {
	ID             string `json:"_key,omitempty"` // increment unique key
	ProviderID     string `json:"providerId"`     // document reference key
	ScheduleID     string `json:"scheduleId"`     // document reference key
	ScheduleRuleID string `json:"scheduleRuleId"` // document reference key
	Title          string `json:"title"`          // 課程標題
	Content        string `json:"content"`        // 課程內容
	// Todo
	// Tag ...
}

type Provider struct {
	ID                  string `json:"_key,omitempty"`      // increment unique key
	LineUserID          string `json:"lineUserId"`          // uniq Line userID (e.g. U123xxxxxxxxdef)
	RealName            string `json:"realName"`            // 中文姓名(真實)
	LineAtName          string `json:"lineAtName"`          // 申請Line官方帳號名稱
	LineAtID            string `json:"lineAtId"`            // 申請Line官方帳號ID(Todo 暫不實作)
	LineAtChannelID     string `json:"lineAtChannelId"`     // 申請Line官方帳號 ChannelID
	LineAtChannelSecret string `json:"lineAtChannelSecret"` // 申請Line官方帳號 ChannelSecret
	LineAtAccessToken   string `json:"lineAtAccessToken"`   // 申請Line官方帳號 AccessToken
	CountryCode         string `json:"countryCode"`         // 手機國碼(Todo 暫不實作)
	LineID              string `json:"lineId"`              // 聯絡的個人Line ID
	PhoneNum            string `json:"phoneNum"`            // 手機號碼
	ConfirmedPhoneNum   string `json:"confirmedPhoneNum"`   // 認證過的手機號碼(Todo 二階段手機驗證)
	GmailAddr           string `json:"gmailAddr"`           // Gamil
	ConfirmedGmailAddr  string `json:"confirmedGmailAddr"`  // 認證過的Gamil(Todo 發認證信)
	GCalSync            bool   `json:"gCalSync"`            // Google Calendar 授權成功
	GUUID               string `json:"guuid"`               // google uuid
	GToken              string `json:"gToken"`              // google token
	GRawData            string `json:"gRawData"`            // google raw data
	InviteCode          string `json:"inviteCode"`          // 企業用戶，業務推廣碼
	InviterID           string `json:"inviterId"`           // invite code對應的provider
	SharableCode        string `json:"sharableCode"`        // 自己的推廣碼
	MemeberTerm         bool   `json:"memeberTerm"`         // 會員條款
	PrivacyTerm         bool   `json:"privacyTerm"`         // 隱私全條款
	Status              int    `json:"status"`              // 狀態 0: 暫存, 1: 確認送出, 2: 審核中, 3: 審核完成, -1: 審核不通過, -2: 例外處理
	CreatedAt           int64  `json:"createdAt"`           // 創建時間
	Blocked             bool   `json:"bolcked"`             // Todo 停權blocked, 停權處理流程
	// Todo 應該還有個google calendar 授權成功拿到的token
}

type Consumer struct {
	ID                      string `json:"_key,omitempty"`          // increment unique key
	ProviderID              string `json:"providerId"`              // document reference key
	LineUserID              string `json:"lineUserId"`              // Line UserId
	ProviderLineAtChannelID string `json:"providerLineAtChannelId"` // Provider Line官方帳號 ChannelID
	LineFollowingStatus     int    `json:"lineFollowingStatus"`     // Line following status
	CreatedAt               int64  `json:"createdAt"`               // 創建時間
}

type Feedback struct {
	ID            string `json:"_key,omitempty"` // increment unique key
	ScheduleID    string `json:"scheduleId"`     // document reference key
	AppointmentID string `json:"appointmentId"`  // document reference key
	Title         string `json:"title"`          // 回饋標題
	Content       string `json:"content"`        // 回饋內容
}

type MonthReceipt struct {
	ID                     string `json:"_key,omitempty"`         // increment unique key
	ProviderID             string `json:"providerId"`             // document reference key
	RawTotalIncome         int32  `json:"rawTotalIncome"`         // 原始收入總額
	RemittedAmount         int32  `json:"remittedAmount"`         // 退款金額
	PlatformRate           int32  `json:"platformRate"`           // 抽成比例
	PlatformFee            int32  `json:"platformFee"`            // 抽成總額
	MarketingFee           int32  `json:"marketingFee"`           // Todo 相關行銷分潤總和
	TotalIncomeGatewayFee  int32  `json:"totalIncomeGatewayFee"`  // 第三方金流收款時所抽取之總額
	TotalOutcomeGatewayFee int32  `json:"totalOutcomeGatewayFee"` // 第三方金流放款時所抽取之總額
	TaxFee                 int32  `json:"taxFee"`                 // Todo 扣掉相關國家稅務
	NetIncome              int32  `json:"netIncome"`              // 老師到帳淨收入
	Notes                  string `json:"notes"`                  // 相關備註
	Paid                   bool   `json:"paid"`                   // 是否放完款項
	ClearingStartedAt      int64  `json:"clearingStartedAt"`      // 結算起始時間
	CreatedAt              int64  `json:"createdAt"`              // 創建時間
	UpdatedAt              int64  `json:"updatedAt"`              // 更新時間
	PaidAt                 int64  `json:"paidAt"`                 // 放款時間
}

type Order struct {
	ID               string `json:"_key,omitempty"`   // increment unique key
	ConsumerID       string `json:"consumerId"`       // document reference key 誰買的
	ProviderID       string `json:"providerId"`       // document reference key 買誰提供的服務
	ServiceProductID string `json:"serviceProductId"` // document reference key 買哪款服務產品
	PaymentMethodID  string `json:"paymentMethodId"`  // document reference key 選擇的結帳方式
	Amount           int    `json:"amount"`           // 購買數量
	Status           int    `json:"status"`           // 狀態, 0初始、1付款中、2已付款、3取消
	CreatedAt        int64  `json:"createdAt"`        // 訂單創建時間
	UpdatedAt        int64  `json:"updatedAt"`        // 更新時間
	PaidAt           int64  `json:"paidAt"`           // 付款日期
}

type Payment struct {
	ID              string `json:"_key,omitempty"`  // increment unique key
	ProviderID      string `json:"providerId"`      // document reference key
	OrderID         string `json:"orderId"`         // document reference key
	ConsumerID      string `json:"consumerId"`      // document reference key
	PaymentMethodID string `json:"paymentMethodId"` // document reference key
	PaidPrice       int32  `json:"paidPrice"`       // 付款金額
	PlatformFee     int32  `json:"platformFee"`     // 我們平台所抽的金額
	PaymentFee      int32  `json:"paymentFee"`      // 金流服務抽成
	AgentFee        int32  `json:"agentFee"`        // Todo 合作抽成(與廠商合作的分潤)
	AdFee           int32  `json:"adFee"`           // Todo 業務推廣抽成
	TaxFee          int32  `json:"taxFee"`          // Todo 勞務報酬報稅
	NetAmount       int32  `json:"netAmount"`       // 可被發放的金額
	Status          int    `json:"status"`          // 狀態, 是否成功付款, 0 init, 1 process, 2 success, -1 異常
	RawParams       string `json:"rawParams"`       // 原始參數
	CreatedAt       int64  `json:"createdAt"`       // 訂單創建時間
	UpdatedAt       int64  `json:"updatedAt"`       // 更新時間
}

// 第三方金流費率
type PaymentMethod struct {
	ID              string `json:"_key,omitempty"`  // increment unique key
	PaymentType     string `json:"paymentType"`     // 付款方式, 刷卡、超商、ATM轉帳...
	ServicePlatform string `json:"servicePlatform"` // 金流平台, default 藍新
	// Todo 看藍新、我們所要支援的種類還有相應的費率
}

type Schedule struct {
	ID               string `json:"_key,omitempty"`   // increment unique key
	ScheduleRuleID   string `json:"scheduleRuleId"`   // document reference key
	CourseID         string `json:"courseId"`         // document reference key
	ProviderID       string `json:"providerId"`       // document reference key
	CourseStartAt    int64  `json:"courseStartAt"`    // 課程開始時間
	CourseEndAt      int64  `json:"courseEndAt"`      // 課程結束時間
	MinConsumerLimit int    `json:"minConsumerLimit"` // 最小開課人數下限
	MaxConsumerLimit int    `json:"maxConsumerLimit"` // 最大開課人數上限
	Count            int    `json:"count"`            // 目前參加人數
	MeetingUrl       string `json:"meetingUrl"`       // 視訊平台連結
}

// CycleRepeatedAmount
// 如果是間隔每天 -> 紀錄0
// 如果是間隔每週 -> 紀錄每個星期幾(1 - 7)禮拜一到禮拜日 -> 先實作這個就好
// 如果是間隔每月 -> 紀錄1 - 31 (每個月的第幾天)
// 如果是間隔每年 -> 紀錄0
// e.g. 每週禮拜五 7/2 10 - 11 am 開出一堂預約 開到11月3號結束
// courseStartAt : "10:00"
// courseEndAt : "11:00"
// cycleStartAt : timestamp 7/2 10:00
// cycleEndAt : timestamp 11/3 10:00
// cycleEndType : -1
// cycleRepeatedAmount : 5
// cycleDiffAmount : 1
// cycleDiffUnit : 1
type ScheduleRule struct {
	ID                  string `json:"_key,omitempty"`      // increment unique key
	ProviderID          string `json:"providerId"`          // document reference key
	CourseStartAt       int64  `json:"courseStartAt"`       // 課程開始時間 // Todo string HH:MM
	CourseEndAt         int64  `json:"courseEndAt"`         // 課程結束時間 // Todo string HH:MM
	CycleStartAt        int64  `json:"cycleStartAt"`        // 週期開始時間
	CycleEndAt          int64  `json:"cycleEndAt"`          // 週期結束時間 -> 配合結束類型看
	CycleRepeatedAmount int    `json:"cycleRepeatedAmount"` // 搭配週期間隔數量、單位看 3(每個禮拜三)
	CycleDiffAmount     int    `json:"cycleDiffAmount"`     // 週期間隔數量 1
	CycleDiffUnit       int    `json:"cycleDiffUnit"`       // 週期間隔單位(0-天、1-週、2-月、3-年)
	CycleEndType        int    `json:"cycleEndType"`        // 結束類型：1幾次、-1時間、-2永遠不停
	MinConsumerLimit    int    `json:"minConsumerLimit"`    // 最小開課人數下限
	MaxConsumerLimit    int    `json:"maxConsumerLimit"`    // 最大開課人數上限
	Count               int    `json:"count"`               // 目前參加人數
}

type ServiceProduct struct {
	ID              string `json:"_key,omitempty"`  // increment unique key
	ProviderID      string `json:"providerId"`      // document reference key
	CountPerPack    int    `json:"countPerPack"`    // 一包多少堂
	Price           int    `json:"price"`           // 一堂多少價格
	ExpiredDuration int64  `json:"expiredDuration"` // 多久過期
	CreatedAt       int64  `json:"createdAt"`       //
	DeletedAt       int64  `json:"deletedAt"`       // soft delete (讓order history可以reference)
}

type ZoomAccount struct {
	UserID    string `json:"userId"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}
