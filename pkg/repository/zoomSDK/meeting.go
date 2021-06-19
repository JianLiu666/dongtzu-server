package zoomSDK

import (
	"dongtzu/config"
	"dongtzu/constant"
	"time"

	"github.com/himalayan-institute/zoom-lib-golang"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetMeetingUrl(startTime time.Time, minute int) (string, int) {
	if !initialized {
		return "", constant.Module_Initialization_Notyet
	}

	handler := getHandler()

	resp, err := handler.Client.CreateMeeting(zoom.CreateMeetingOptions{
		HostID:    handler.UserId,
		Topic:     "DongTzu Meeting",
		Type:      zoom.MeetingTypeScheduled,
		StartTime: &zoom.Time{Time: startTime},
		Duration:  minute + config.GetGlobalConfig().Zoom.MeetingExtendedTime,
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
		return "", constant.ZoomSDK_Driver_Failed
	}

	return resp.JoinURL, constant.ZoomSDK_Success
}
