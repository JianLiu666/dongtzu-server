package config

type ZoomConfig struct {
	BaseUrl   string `yaml:"baseUrl,omitempty"`
	UserID    string `yaml:"userId,omitempty"`
	ApiKey    string `yaml:"apiKey,omitempty"`
	ApiSecret string `yaml:"apiSecret,omitempty"`
	JwtToken  string `yaml:"jwtToken,omitempty"`
}
