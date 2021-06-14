# ArangoDB Schema

## Providers
| Field          | Type   | Description |
| -------------- | ------ | ----------- |
| _key           | string | increment unique key |
| status         | string | 0: 暫存, 1: 確認送出, 2: 審核中, 3: 審核完成, 4: 審核不通過, 5: 例外處理 |