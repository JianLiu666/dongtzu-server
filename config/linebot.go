package config

type lineBotConfig struct {
	ChannelSecret      string `yaml:"channelSecret,omitempty"`
	ChannelAccessToken string `yaml:"channelAccessToken,omitempty"`
	Port               string `yaml:"port,omitempty"`
}
