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
- [ ] 向 zoom 申請一個 meeting url
- [x] 批次更新 schedules 到 db

**對準備開始的 appt 發送 meeting url**

- [ ] 從 db 撈取準備開始的 appts (用 schedules 關聯查詢)
- [ ] 發送連結到 line 用戶 (串接 push message)
- [ ] 批次更新以發過連結的 appts 的 status 到 db (確保不用在重新發送, 節省流量)

**對發送準備結束的 appt 發送 feedback url**

- [ ] 從 db 撈取準備結束的 appts (用 schedules 關聯查詢)
- [ ] 發送連結到 line 用戶 (由於這個 webview url 應該不會變, 開一個 collection 紀錄這個連結)
- [ ] 批次更新已發過連結的 appts 的 status 到 db (確保不用在重新發送, 節省流量)