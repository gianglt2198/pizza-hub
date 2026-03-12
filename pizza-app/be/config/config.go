package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type AppConfig struct {
	// Server
	ServerPort string `json:"server_port" mapstructure:"server_port"`
	Env        string `json:"env" mapstructure:"dev"`
}

type DatabaseConfig struct {
	// PostgreSQL
	DBHost     string `mapstructure:"host"`
	DBPort     string `mapstructure:"port"`
	DBUser     string `mapstructure:"user"`
	DBPassword string `mapstructure:"password"`
	DBName     string `mapstructure:"name"`
	DBSSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	// Redis
	RedisHost     string `mapstructure:"host"`
	RedisPort     string `mapstructure:"port"`
	RedisPassword string `mapstructure:"password"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	configPath := os.Getenv("CONFIG_PATH")

	// Set config file path
	if configPath != "" {
		v.AddConfigPath(configPath)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Enable environment variable support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDatabaseURL returns a formatted database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		"postgres",
		c.Database.DBUser,
		c.Database.DBPassword,
		c.Database.DBHost,
		c.Database.DBPort,
		c.Database.DBName,
		c.Database.DBSSLMode,
	)
}

func (c *Config) RedisAddr() string {
	return c.Redis.RedisHost + ":" + c.Redis.RedisPort
}
