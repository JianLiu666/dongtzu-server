package config

type ZoomConfig struct {
	ApiKey    string `yaml:"apiKey,omitempty"`
	ApiSecret string `yaml:"apiSecret,omitempty"`
	JwtToken  string `yaml:"jwtToken,omitempty"`
}
