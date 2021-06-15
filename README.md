# dongtzu-server

- [dongtzu-server](#dongtzu-server)
  - [TODO](#todo)
    - [Jian6](#jian6)
      - [Scheduler](#scheduler)

----

## TODO

### Jian6

#### Scheduler

**定時撈取近期尚未建立 zoom meeting url 的 schedule**

- [x] 從 db 撈取未處理的 schedules
- [X] 向 zoom 申請一個 meeting url
- [x] 批次更新 schedules 到 db

**對準備開始的 appt 發送 meeting url**

- [x] 從 db 撈取準備開始的 appts (用 schedules 關聯查詢)
- [x] 發送連結到 line 用戶 (串接 push message)
- [ ] message type 應該要改為圖文訊息(Image Message) 效果會更好
- [x] 批次更新以發過連結的 appts 的 status 到 db (確保不用在重新發送, 節省流量)

**對發送準備結束的 appt 發送 feedback url**

- [x] 從 db 撈取準備結束的 appts (用 schedules 關聯查詢)
- [x] 發送 mock 訊息
- [ ] 發送連結到 line 用戶 (由於這個 webview url 應該不會變, 開一個 collection 紀錄這個連結)
- [ ] message type 應該要改為圖文訊息(Image Message) 效果會更好
- [x] 批次更新已發過連結的 appts 的 status 到 db (確保不用在重新發送, 節省流量)

**串接 ZoomSDK Webhook**

之後可以詳閱一下 webhook event, 即時監聽正在進行中的 meeting 狀況

**監聽追蹤老師帳號的LineID**

- [ ] 在學生追蹤老師的當下即刻在 Consumers 建立一筆 doc (後續在付費or正式成為會員時再繼續更新這個 doc)

**驗證 Schedule 合法性**

- [ ] 確保 consumer 再預約課程時不會跟自己的其他課程相撞