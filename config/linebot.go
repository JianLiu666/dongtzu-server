package config

type LineBotConfig struct {
	ChannelSecret      string `yaml:"channelSecret,omitempty"`
	ChannelAccessToken string `yaml:"channelAccessToken,omitempty"`
}
