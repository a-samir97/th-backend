package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Elasticsearch ElasticsearchConfig
	Redis        RedisConfig
	Queue        QueueConfig
	Storage      StorageConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ElasticsearchConfig struct {
	URL   string
	Index string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type QueueConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type StorageConfig struct {
	Type      string // "local" or "s3"
	LocalPath string
	S3Bucket  string
	S3Region  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "thamaniyah"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Elasticsearch: ElasticsearchConfig{
			URL:   getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
			Index: getEnv("ELASTICSEARCH_INDEX", "media"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Queue: QueueConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
			User:     getEnv("RABBITMQ_USER", "admin"),
			Password: getEnv("RABBITMQ_PASSWORD", "admin"),
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "local"),
			LocalPath: getEnv("STORAGE_LOCAL_PATH", "./uploads"),
			S3Bucket:  getEnv("STORAGE_S3_BUCKET", ""),
			S3Region:  getEnv("STORAGE_S3_REGION", "us-east-1"),
		},
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) QueueURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.Queue.User,
		c.Queue.Password,
		c.Queue.Host,
		c.Queue.Port,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
