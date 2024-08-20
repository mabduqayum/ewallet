package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Wallet   WalletConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port    int
	Host    string
	Version string
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
	SSLMode  string
	Schema   string
}

func (d DatabaseConfig) ConnectionString() string {
	query := url.Values{}
	query.Add("sslmode", d.SSLMode)
	query.Add("search_path", d.Schema)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:     d.DBName,
		RawQuery: query.Encode(),
	}

	return u.String()
}

type WalletConfig struct {
	UnidentifiedLimit float64
	IdentifiedLimit   float64
}

type LoggingConfig struct {
	Level string
}

func LoadConfig(env string) (*Config, error) {
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath("./internal/config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Override with environment variables if they exist
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}

	return &config, nil
}
