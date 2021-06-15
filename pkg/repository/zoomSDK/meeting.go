package zoomSDK

import (
	"dongtzu/config"
	"dongtzu/constant"
	"time"

	"github.com/himalayan-institute/zoom-lib-golang"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetMeetingUrl(startTime time.Time, minute int) (string, int) {
	resp, err := client.CreateMeeting(zoom.CreateMeetingOptions{
		HostID:    config.GetGlobalConfig().Zoom.UserID,
		Topic:     "DongTzu Meeting",
		Type:      zoom.MeetingTypeScheduled,
		StartTime: &zoom.Time{Time: startTime},
		Duration:  minute,
		Timezone:  "Asia/Taipei",
		Password:  "",
		Agenda:    "Test meeeting api from zoom sdk.",
		Settings: zoom.MeetingSettings{
			ParticipantVideo: true, // start video when participants join the meeting.
			JoinBeforeHost:   true, // allow participants to join the meeting before the host starts the meeting.
		},
	})

	if err != nil {
		logger.Errorf("[ZoomSDK] create meeting url failed: %v", err)
		return "", constant.Zoom_Driver_Failed
	}

	return resp.JoinURL, constant.Zoom_Success
}
