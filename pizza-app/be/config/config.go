package config

import (
	"fmt"
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
	DBHost     string `envconfig:"db_host" default:"localhost"`
	DBPort     string `envconfig:"db_port" default:"5432"`
	DBUser     string `envconfig:"db_user" default:"pizza"`
	DBPassword string `envconfig:"db_password" default:"pizza123"`
	DBName     string `envconfig:"db_name" default:"pizzahub"`
	DBSSLMode  string `envconfig:"db_ssl_mode" default:"disable"`
}

type RedisConfig struct {
	// Redis
	RedisHost     string `envconfig:"redis_host" default:"localhost"`
	RedisPort     string `envconfig:"redis_port" default:"6379"`
	RedisPassword string `envconfig:"redis_password" default:""`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

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
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
