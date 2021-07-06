# Repository

- [Repository](#repository)
  - [Summary](#summary)
  - [TODOs](#todos)
    - [Basic](#basic)
    - [Database Migration](#database-migration)
    - [Zoom SDK](#zoom-sdk)
    - [Line SDK](#line-sdk)
    - [NewebPay SDK](#newebpay-sdk)

<br>

## Summary

施工中 ...

<br>

## TODOs

### Basic

- [ ] 使 Viper 能夠動態更新參數

### Database Migration

- [x] AranogoDB Seed
- [ ] 單元測試環境

### Zoom SDK

- [ ] 購買 Zoom Meeting Account PRO Plan : https://zoom.us/pricing
- [x] 啟動時從資料庫取得 Zoom Host Account
  - [x] 單一帳號 (串測)
  - [x] 管理多個帳號 (Account Pool)
    - [ ] 紀錄每個帳號在當天的累計請求數量 (00:00 - 23:59)
    - [x] Round-robin 使用每個帳號發送 API
    - [ ] 提供學員免費試用功能 (Free Account 單次會議時間上限:40分鐘)
- [x] Create Scheduled Meeting
  - [x] 設定時間區間 (提供 n+10mins 作為課後時間)
  - [x] Host 無需事先進入也可以讓 Consumer 進入
  - [ ] 允許 MeetingURL 可以提前進入 (e.g. 提前10分鐘進入會議室準備)
  - [ ] 測試 response 響應時間 (目前看來開一間房間至少要 1s 的時間)
- [ ] Monitor Starting Meetings
  - [ ] 確認當前課程人數是否符合最低人數限制
  - [ ] 提早結束空堂課程
  - [ ] 儲存會議備份 (供學員課後回放)

### Line SDK

- [x] 啟動時從資料庫取得 Line Provider Accounts
  - [x] 單一帳號 (串測)
  - [x] 管理多個帳號 (Account Pool)
    - [x] 從 webhook subdomain 識別不同 provider 傳來的 request
- [x] Webhook Event
  - [x] OnFollow
    - [x] 至資料庫建立會員(目的是保留UserLineId)
  - [x] OnUnfollow
    - [x] 至資料庫軟刪除(只修改lineFollowingStatus)
- [x] Push Message
  - [x] 發送訊息至指定學員
    - [x] Text Message(串測)
    - [x] Flex Message
      - [x] 串測 API
      - [x] 設定跳轉連結(MeetingURL、FeedbackURL)
- [ ] 從資料庫載入 Flex Template 

### NewebPay SDK

- [ ] 串測(待定)