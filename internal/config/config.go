package config

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type ServerConfig struct {
	Host         string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Port         string        `yaml:"port" env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"1m"`
}

type StorageConfig struct {
	Path string `yaml:"path" env:"STORAGE_PATH" env-default:"./storage"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Logger  LoggerConfig  `yaml:"logger"`
	API     API           `yaml:"api"`
}

func (srv *ServerConfig) ServerAddr() string {
	return net.JoinHostPort(srv.Host, srv.Port)
}

func (c *Config) Format() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Server:\n")
	fmt.Fprintf(&b, "  Host:           %s\n", c.Server.Host)
	fmt.Fprintf(&b, "  Port:           %s\n", c.Server.Port)
	fmt.Fprintf(&b, "  ReadTimeout:    %s\n", c.Server.ReadTimeout)
	fmt.Fprintf(&b, "  WriteTimeout:   %s\n", c.Server.WriteTimeout)
	fmt.Fprintf(&b, "  IdleTimeout:    %s\n", c.Server.IdleTimeout)

	fmt.Fprintf(&b, "\nStorage:\n")
	fmt.Fprintf(&b, "  Path:           %s\n", c.Storage.Path)

	fmt.Fprintf(&b, "\nLogger:\n")
	fmt.Fprintf(&b, "  Level:          %s\n", c.Logger.Level)

	fmt.Fprintf(&b, "\n*****\nAPI:\n")
	fmt.Fprintf(&b, "%s", c.API.Format())

	return b.String()
}

type API struct {
	Task TaskConfig `yaml:"task"`
	Arch ArchConfig `yaml:"arch"`
}

func (c *API) Format() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Task:\n")
	fmt.Fprintf(&b, "  Timeout:        %s\n", c.Task.Timeout)
	fmt.Fprintf(&b, "  MaxTask:        %d\n", c.Task.MaxTask)
	fmt.Fprintf(&b, "  MaxURL:         %d\n", c.Task.MaxURL)
	fmt.Fprintf(&b, "  Ext:            %s\n", c.Task.Ext)
	fmt.Fprintf(&b, "  Arch:           %s\n", c.Task.Arch)
	fmt.Fprintf(&b, "  Path:           %s\n", c.Task.Path)

	fmt.Fprintf(&b, "\nArch:\n")
	fmt.Fprintf(&b, "  Path:           %s\n", c.Arch.Path)

	return b.String()
}

type TaskConfig struct {
	Timeout time.Duration `yaml:"timeout" env:"TASK_TIMEOUT" env-default:"5m"`
	MaxTask int           `yaml:"max_task" env:"TASK_MAX" env-default:"3"`
	MaxURL  int           `yaml:"max_url" env:"URL_MAX" env-default:"3"`
	Ext     []string      `yaml:"ext" env:"EXT"`
	Arch    string        `yaml:"arch" env:"ARCH"`
	Path    string        `yaml:"path" env:"PATH_TASK" env-default:"tmp"`
}

type ArchConfig struct {
	Path string `yaml:"path" env:"PATH_ARCH" env-default:"tmp"`
}
