# APIServer

- [APIServer](#apiserver)
  - [Summary](#summary)
  - [TODOs](#todos)

<br>

## Summary

處理透過 Provider Line 官方帳號傳送來請求：
 - 購買課程
 - 預約課程

<br>

## TODOs

- [ ] [API] 取得老師的授課方案(ServiceProduct)
- [ ] [API] 購買課程(待定)
- [ ] [API] 預約課程
  - [ ] API Request Schema
  - [ ] API Response Schema
  - [ ] 處理流程
    - [ ] 檢查 Payload 是否符合格式
    - [x] 寫入資料庫 (transaction)
      - [x] 1. 檢查指定 Schedule 是否符合預約規定 (e.g. 人數上限)
      - [x] 2. 檢查 Consumer 在相同時段是否存在其他 Appointments
      - [ ] 3. 檢查 Consumer 是否還有足夠的剩餘堂數
      - [x] 4. create Appointment
      - [x] 5. 更新 schedule count(當前人數)
    - [ ] 傳送 appointment 到 scheduler cache (對應 [Scheduler 處理 appointments 的解 3](./scheduler.md#scheduler-1))
    - [ ] Response
