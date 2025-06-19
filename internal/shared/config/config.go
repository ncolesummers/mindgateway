package config

import (
	"fmt"
	"os"
	"time"
	
	"github.com/spf13/viper"
)

type Config struct {
	// General settings
	Environment     string        `mapstructure:"environment"`
	LogLevel        string        `mapstructure:"log_level"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	
	// Server settings
	Server struct {
		Address string `mapstructure:"address"`
		Port    int    `mapstructure:"port"`
	} `mapstructure:"server"`
	
	// Database settings
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"database"`
	
	// Redis settings
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	
	// ETCD settings
	ETCD struct {
		Endpoints []string `mapstructure:"endpoints"`
		Username  string   `mapstructure:"username"`
		Password  string   `mapstructure:"password"`
	} `mapstructure:"etcd"`
	
	// Service specific settings
	Auth struct {
		Address   string `mapstructure:"address"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"auth"`
	
	Registry struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"registry"`
	
	// Worker settings
	Worker struct {
		ConnectTimeout    time.Duration `mapstructure:"connect_timeout"`
		RequestTimeout    time.Duration `mapstructure:"request_timeout"`
		HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
	} `mapstructure:"worker"`
	
	// Queue settings
	Queue struct {
		MaxSize          int           `mapstructure:"max_size"`
		DefaultPriority  int           `mapstructure:"default_priority"`
		ProcessingPeriod time.Duration `mapstructure:"processing_period"`
	} `mapstructure:"queue"`
}

// Load loads the configuration from file and environment
func Load() (*Config, error) {
	cfg := &Config{}
	
	// Set default values
	setDefaults()
	
	// Get config file path from environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Default to configs directory
		configPath = "configs"
	}
	
	// Get environment from ENV var
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev" // Default to dev environment
	}
	
	// Setup viper
	viper.AddConfigPath(configPath)
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	
	// Read environment variables
	viper.AutomaticEnv()
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// General defaults
	viper.SetDefault("environment", "dev")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("shutdown_timeout", 30*time.Second)
	
	// Server defaults
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "mindgateway")
	viper.SetDefault("database.name", "mindgateway")
	viper.SetDefault("database.ssl_mode", "disable")
	
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	
	// ETCD defaults
	viper.SetDefault("etcd.endpoints", []string{"localhost:2379"})
	
	// Service defaults
	viper.SetDefault("auth.address", "localhost:9091")
	viper.SetDefault("registry.address", "localhost:9092")
	
	// Worker defaults
	viper.SetDefault("worker.connect_timeout", 5*time.Second)
	viper.SetDefault("worker.request_timeout", 60*time.Second)
	viper.SetDefault("worker.health_check_period", 30*time.Second)
	
	// Queue defaults
	viper.SetDefault("queue.max_size", 10000)
	viper.SetDefault("queue.default_priority", 5)
	viper.SetDefault("queue.processing_period", 100*time.Millisecond)
}