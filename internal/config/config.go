package config

import "github.com/spf13/viper"

type Config struct {
	RabbitMQ struct {
		URL   string `mapstructure:"url"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"rabbitmq"`

	Email struct {
		SMTPHost       string `mapstructure:"smtp_host"`
		SMTPPort       int    `mapstructure:"smtp_port"`
		SenderEmail    string `mapstructure:"sender_email"`
		SenderPassword string `mapstructure:"sender_password"`
		RecipientEmail string `mapstructure:"recipient_email"`
		Subject        string `mapstructure:"subject"`
	} `mapstructure:"email"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
