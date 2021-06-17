# Scheduler

- [Scheduler](#scheduler)
  - [Summary](#summary)
  - [TODOs](#todos)
    - [Startup](#startup)
    - [Zoom SDK](#zoom-sdk)
    - [Line SDK](#line-sdk)
    - [Scheduler](#scheduler-1)

<br>

## Summary

用來處理週期性任務，例如：

- 對所有滿足上課條件的課程至 Zoom 建立 ScheduledMeeting。
- 對所有準備開始的課程發送 MeetingURL 至對應的學員 Line 帳號。
- 對所有準備結束的課程發送 FeedbackURL 至對應的學員 Line 帳號。

<br>

## TODOs

### Startup

- [ ] 使 Viper 能夠動態更新參數

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

- [ ] 啟動時從資料庫取得 Line Provider Accounts
  - [ ] 單一帳號 (串測)
  - [ ] 管理多個帳號 (Account Pool)
    - [ ] 從 webhook subdomain 識別不同 provider 傳來的 request
- [ ] Webhook Event
  - [ ] OnFollow
    - [ ] 至資料庫建立會員(目的是保留UserLineId)
  - [ ] OnUnfollow
    - [ ] 至資料庫軟刪除(只修改lineFollowingStatus)
- [ ] Push Message
  - [ ] 發送訊息至指定學員
    - [ ] Text Message(串測)
    - [ ] ImageMap Message
      - [ ] 串測 API
      - [ ] 設定跳轉連結(MeetingURL、FeedbackURL)

### Scheduler

- [ ] 循環事件
  - [ ] 替未申請 MeetingURL 的 Schedule 申請 MeetingURL
    - [ ] 只處理即將到來的 Schedule (根據時間參數決定)
    - [ ] 每 n 分鐘處理一次 (根據時間參數、會議建立速度決定)
  - [ ] 會議開始前 n 分鐘發送 MeetingURL 至 consumers 的 line
    - [ ] 每個 20, 50 分鐘時取得準備開始的 appointments 
      - [ ] 取得 appointment 方式
        - [ ] 從資料庫一次撈取所有符合條件的 appointments (解1)
        - [ ] 批次撈取所有符合條件的 appointments (解2, 降低單次資料量)
        - [ ] 預先 cache 即將到來的 appointments 到 timeWheel 上，直到開始前發送 (解3, 需要災難還原機與 appintment 資料被更新是否有影響)
    - [ ] 發送訊息至對應的 consumers line
      - [ ] Text Message (串測)
      - [ ] ImageMap Message
    - [ ] 將 appointment 的狀態改為已發送 MeetingURL (確保不會重複發送)
    - [ ] 將已處理的 appointments 更新回資料庫
  - [ ] 會議結束前 n 分鐘發送 FeedbackURL 至 consumers 的 line
    - [ ] 每個 20, 50 分鐘時取得準備結束的 appointments 
      - [ ] 取得 appointment 方式
        - [ ] 從資料庫一次撈取所有符合條件的 appointments (解1)
        - [ ] 批次撈取所有符合條件的 appointments (解2, 降低單次資料量)
        - [ ] 預先 cache 即將到來的 appointments 到 timeWheel 上，直到開始前發送 (解3, 需要災難還原機與 appintment 資料被更新是否有影響)
    - [ ] 發送訊息至對應的 consumers line
      - [ ] Text Message (串測)
      - [ ] ImageMap Message
    - [ ] 將 appointment 的狀態改為已發送 FeedbackURL (確保不會重複發送)
    - [ ] 將已處理的 appointments 更新回資料庫
