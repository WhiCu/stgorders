package config

import (
	"net"
	"strings"
)

type ServerConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
}

type StorageConfig struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User string `yaml:"user" env:"DB_USER" env-default:"user"`
	Pass string `yaml:"pass" env:"DB_PASS" env-default:"password"`
	Name string `yaml:"name" env:"DB_NAME" env-default:"l0_test"`
}

func (c *StorageConfig) DSN() string {
	return "postgresql://" + c.User + ":" + c.Pass + "@" + c.Host + ":" + c.Port + "/" + c.Name + "?sslmode=disable"
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Path  string `yaml:"path" env:"LOG_PATH" env-default:""`
	Size  int    `yaml:"size" env:"LOG_FILE_SIZE" env-default:"128"`
}

type WorkerPoolConfig struct {
	Size int `yaml:"size" env:"WORKER_POOL_SIZE" env-default:"128"`
	Buf  int `yaml:"buf" env:"WORKER_POOL_BUF" env-default:"128"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"test"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"test"`

	WorkerPool WorkerPoolConfig `yaml:"worker_pool"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Logger  LoggerConfig  `yaml:"logger"`
	Kafka   KafkaConfig   `yaml:"kafka"`
}

func (srv *ServerConfig) ServerAddr() string {
	return net.JoinHostPort(srv.Host, srv.Port)
}

func (c *Config) Format() string {
	var b strings.Builder

	return b.String()
}
