# dongtzu-server

- [dongtzu-server](#dongtzu-server)
  - [Introduction](#introduction)
  - [Key Features](#key-features)
  - [Prototype](#prototype)
  - [Server Setup](#server-setup)

<br>

## Introduction

台灣在2021年5月28日因為疫情爆發的關係，宣布正式進入全國第三級防疫警戒，造成多數行業頓時失去經濟收入。  
為了使健身教練能夠在家中也能實現遠端授課，且無需額外安裝任何多餘的APP，快速達到其需求目的。

因此，整合下列平台與第三方工具，成為一個線上預約系統的後端服務：

**Core**

- API Server
- Scheduler

**3rd-party Tool**

- Line Webhook SDK
- NewebPay SDK
- Zoom SDK

<br>

## Key Features

- 透過 Line 官方帳號，提供教練與學生付費、預約與上課管道。
- 透過 NewebPay 即時支付。
- 透過 Zoom Meeting 功能提供教練與學生視訊上課空間。

<br>

## Prototype

**DongTzu Official Account**

![DongTzu](./doc/img/DongTzu_OfficialAccount.svg)

**Teacher Official Account**

![Teacher](./doc/img/Teacher_OfficialAccount.svg)

<br>

## Server Setup

```shell
go run main.go server
```