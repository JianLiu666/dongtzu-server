# APIServer

- [APIServer](#apiserver)
  - [Summary](#summary)
  - [TODOs](#todos)
    - [NewebPay SDK](#newebpay-sdk)
    - [APIServer](#apiserver-1)
    - [Migration](#migration)

<br>

## Summary

處理透過 Provider Line 官方帳號傳送來請求：
 - 購買課程
 - 預約課程

<br>

## TODOs

### NewebPay SDK

- [ ] 串測(待定)

### APIServer

- [ ] [API] 購買課程(待定)
- [ ] [API] 預約課程
  - [ ] API Request Schema
  - [ ] API Response Schema
  - [ ] 處理流程
    - [ ] 檢查 Payload 是否符合格式
    - [ ] 寫入資料庫 (transaction)
      - [ ] 1. 檢查指定 Schedule 是否符合預約規定 (e.g. 人數上限)
      - [ ] 2. 檢查 Consumer 在相同時段是否存在其他 Appointments
      - [ ] 3. create Appointment
      - [ ] 4. 更新 schedule count(當前人數)
    - [ ] 傳送 appointment 到 scheduler cache (對應 [Scheduler 處理 appointments 的解 3](./scheduler.md#scheduler-1))
    - [ ] Response

### Migration

- [ ] AranogoDB Seed
- [ ] 單元測試環境

### 藍新金流測試帳密
cwww.newebay.com
- ts0542648@gmail.com
- P@ssw0rd