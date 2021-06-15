package config

type ZoomConfig struct {
	UserID    string `yaml:"userId,omitempty"`
	ApiKey    string `yaml:"apiKey,omitempty"`
	ApiSecret string `yaml:"apiSecret,omitempty"`
	JwtToken  string `yaml:"jwtToken,omitempty"`
}
