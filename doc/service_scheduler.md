# Scheduler

- [Scheduler](#scheduler)
  - [Summary](#summary)
  - [TODOs](#todos)

<br>

## Summary

用來處理週期性任務，例如：

- 對所有滿足上課條件的課程至 Zoom 建立 ScheduledMeeting。
- 對所有準備開始的課程發送 MeetingURL 至對應的學員 Line 帳號。
- 對所有準備結束的課程發送 FeedbackURL 至對應的學員 Line 帳號。****

<br>

## TODOs

- [x] 循環事件
  - [x] 替未申請 MeetingURL 的 Schedule 申請 MeetingURL
    - [x] 只處理即將到來的 Schedule (根據時間參數決定)
    - [x] 每隔 n 分鐘處理一次 (根據時間參數、會議建立速度決定)
  - [x] 會議開始前 n 分鐘發送 MeetingURL 至 consumers 的 line
    - [x] 每隔 n 分鐘時取得準備開始的 appointments 
      - [x] 取得 appointment 方式
        - [x] 從資料庫一次撈取所有符合條件的 appointments (解1)
        - [ ] 批次撈取所有符合條件的 appointments (解2, 降低單次資料量)
        - [ ] 預先 cache 即將到來的 appointments 到 timeWheel 上，直到開始前發送 (解3, 需要災難還原機與 appintment 資料被更新是否有影響)
    - [x] 發送訊息至對應的 consumers line
      - [x] Text Message (串測)
      - [x] Flex Message
    - [x] 將 appointment 的狀態改為已發送 MeetingURL (確保不會重複發送)
    - [x] 將已處理的 appointments 更新回資料庫
