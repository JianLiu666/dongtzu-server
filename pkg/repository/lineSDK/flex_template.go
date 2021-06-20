package lineSDK

import "fmt"

func getMeetingFlexTemplate(meetingUrl string) string {
	jsonString := `{
		"type": "bubble",
		"size": "giga",
		"hero": {
		  "type": "image",
		  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover"
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "Course Title",
			  "weight": "bold",
			  "size": "xl"
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "margin": "md",
			  "contents": [
				{
				  "type": "text",
				  "text": "About course description.",
				  "size": "sm",
				  "color": "#999999",
				  "margin": "md",
				  "flex": 0
				}
			  ]
			}
		  ]
		},
		"footer": {
		  "type": "box",
		  "layout": "horizontal",
		  "spacing": "sm",
		  "contents": [
			{
			  "type": "button",
			  "style": "primary",
			  "height": "md",
			  "action": {
				"type": "uri",
				"label": "開始上課",
				"uri": "%s"
			  },
			  "color": "#333333"
			},
			{
			  "type": "button",
			  "style": "primary",
			  "height": "md",
			  "action": {
				"type": "postback",
				"label": "結束簽退",
				"data": "GET_FEEDBACK_URL"
			  },
			  "color": "#333333"
			}
		  ],
		  "flex": 0
		},
		"direction": "ltr"
	}`

	return fmt.Sprintf(jsonString, meetingUrl)
}

func getFeedbackFlexTemplate(feedbackUrl string) string {
	jsonString := `{
		"type": "bubble",
		"size": "giga",
		"hero": {
		  "type": "image",
		  "url": "https://memegenerator.net/img/instances/60468871.jpg",
		  "size": "full",
		  "aspectRatio": "1:1",
		  "aspectMode": "cover",
		  "action": {
			"type": "uri",
			"label": "action",
			"uri": "%v"
		  }
		},
		"direction": "ltr"
	}`

	return fmt.Sprintf(jsonString, feedbackUrl)
}
