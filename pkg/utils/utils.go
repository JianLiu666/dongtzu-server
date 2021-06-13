package utils

import "time"

// 已當前時間取得時間區間(UTC)
func GetTimeRange() (int64, int64) {
	var startTimestamp, endTimestamp int64
	now := time.Now().UTC()

	if now.Minute() < 30 {
		now = now.Round(time.Hour)
		startTimestamp = now.Add(30 * time.Minute).UTC().Unix()
		endTimestamp = now.Add(59 * time.Minute).UTC().Unix()
	} else {
		now = now.Round(time.Hour)
		startTimestamp = now.UTC().Unix()
		endTimestamp = now.Add(29 * time.Minute).UTC().Unix()
	}

	return startTimestamp, endTimestamp
}
