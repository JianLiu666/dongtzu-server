package utils

import "time"

// 已當前時間取得時間區間
func GetTimeRange() (int64, int64) {
	var startTimestamp, endTimestamp int64
	now := time.Now()

	if now.Minute() < 30 {
		now = now.Round(time.Hour)
		startTimestamp = now.Add(30 * time.Minute).Unix()
		endTimestamp = now.Add(59 * time.Minute).Unix()
	} else {
		now = now.Round(time.Hour)
		startTimestamp = now.Unix()
		endTimestamp = now.Add(29 * time.Minute).Unix()
	}

	return startTimestamp, endTimestamp
}

// 無條件進位到下個最近的半點時間
//
// @return int64 timestmap
func GetTimestampRoundToNextHalf() int64 {
	now := time.Now()
	if now.Minute() < 30 {
		return now.Round(time.Hour).Add(30 * time.Minute).Unix()
	}
	return now.Round(time.Hour).Unix()
}
