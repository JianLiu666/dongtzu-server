# Feature

- [Feature](#feature)
  - [測試購買](#測試購買)
  - [測試預約](#測試預約)
  - [測試流程](#測試流程)

<br>

## 測試購買

- [x] 發送文字訊息，由 Webhook 轉拋至 Server
- [x] 查詢 `Consumers`，對不存在的 LineUserId 建立一個專屬的 Document
- [ ] 預建 `ServiceProduct` 與 `PaymentMethod` mock data
- [ ] 跳過第三方金流串接，實作建立 `Order` 與 `Payment` 流程

<br>

## 測試預約

- [ ] 發送文字訊息，由 Webhook 轉拋至 Server
- [ ] 透過 Migration 預建時間範圍內每隔30分鐘一筆的 Schedule mock data
- [ ] 查詢 `Consumers`，對不存在的 LineUserId 建立一個專屬的 Document
- [ ] 立即預約距離最近的 Schedule，實作建立 `Appointment` 流程

<br>

## 測試流程

- [ ] 發送文字訊息，由 Webhook 轉拋至 Server
- [ ] 主動發送一筆 FlexMessage 至對應的 UserId