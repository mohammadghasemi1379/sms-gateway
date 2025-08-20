package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Redis    RedisConfig
	App      AppConfig
	Log      LogConfig
	Database DatabaseConfig
	Mock     MockConfig
	RabbitMQ RabbitMQConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type AppConfig struct {
	Environment string
	Port        string
	Name        string
	Version     string
	Host        string
}

type LogConfig struct {
	Level  string
	Format string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type MockConfig struct {
	Host string
	Port int
}

type RabbitMQConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	VHost         string
	PrefetchCount int
}

type ThrottleConfig struct {
	MaxMessagesPerSecond int
	QueueName            string
	PriorityQueues       map[string]int // queue name -> weight percentage
}

func Load() *Config {
	godotenv.Load(".env")
	return &Config{
		Redis:    loadRedisConfig(),
		App:      loadAppConfig(),
		Log:      loadLogConfig(),
		Database: loadDatabaseConfig(),
		Mock:     loadMockConfig(),
		RabbitMQ: loadRabbitMQConfig(),
	}
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		Host:        getEnv("APP_HOST", "0.0.0.0"),
		Name:        getEnv("APP_NAME", "sms-gateway"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
	}
}

func loadLogConfig() LogConfig {
	return LogConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 3306),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "sms_gateway"),
	}
}

func loadMockConfig() MockConfig {
	return MockConfig{
		Host: getEnv("MOCK_HOST", "localhost"),
		Port: getEnvAsInt("MOCK_PORT", 8081),
	}
}

func loadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		Host:     getEnv("RABBITMQ_HOST", "localhost"),
		Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
		User:     getEnv("RABBITMQ_USER", "guest"),
		Password: getEnv("RABBITMQ_PASSWORD", "guest"),
		VHost:    getEnv("RABBITMQ_VHOST", "/"),
		PrefetchCount: getEnvAsInt("RABBITMQ_PREFETCH_COUNT", 10),
	}
}

func (r *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (a *AppConfig) IsProduction() bool {
	return strings.ToLower(a.Environment) == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
