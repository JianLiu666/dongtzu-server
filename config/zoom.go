package config

type ZoomConfig struct {
	MeetingExtendedTime int    `yaml:"meetingExtendedTime"`
	DefaultUserID       string `yaml:"defaultUserId"`
	DefaultAPIKey       string `yaml:"defaultAPIKey"`
	DefaultAPISecret    string `yaml:"defaultAPISecret"`
}
