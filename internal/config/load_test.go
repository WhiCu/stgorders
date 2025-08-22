package config

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const sampleYAML = `
server:
  host: "127.0.0.1"
  port: "9090"
  read_timeout: 1s
  write_timeout: 2s
  idle_timeout: 3s
storage:
  host: "db"
  port: "5432"
  user: "u"
  pass: "p"
  name: "n"
logger:
  level: "debug"
  path: ""
  size: 128
  compress: false
kafka:
  brokers: ["k1:9092","k2:9092"]
  topic: "t"
  group_id: "g"
  worker_pool:
    size: 2
    buf: 4
cache:
  size: 16
`

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestLoad(t *testing.T) {
	Convey("LoadWithDefault loads config from provided path when env empty", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Setenv(configPath, path)
		defer os.Unsetenv(configPath)
		cfg, err := Load()
		So(err, ShouldBeNil)
		So(cfg.Server.Port, ShouldEqual, "9090")
		So(cfg.Storage.User, ShouldEqual, "u")
		So(cfg.Kafka.WorkerPool.Size, ShouldEqual, 2)
	})
}

func TestLoadWithDefault(t *testing.T) {
	Convey("LoadWithDefault loads config from provided path when env empty", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Unsetenv(configPath)
		cfg, err := LoadWithDefault(path)
		So(err, ShouldBeNil)
		So(cfg.Server.Port, ShouldEqual, "9090")
		So(cfg.Storage.User, ShouldEqual, "u")
		So(cfg.Kafka.WorkerPool.Size, ShouldEqual, 2)
	})
}

func TestLoadLogger(t *testing.T) {
	Convey("LoadLogger returns logger subconfig", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Setenv(configPath, path)
		defer os.Unsetenv(configPath)
		lc, err := LoadLogger()
		So(err, ShouldBeNil)
		So(lc.Level, ShouldEqual, "debug")
	})
}

func TestLoadServer(t *testing.T) {
	Convey("LoadServer returns server subconfig", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Setenv(configPath, path)
		defer os.Unsetenv(configPath)
		sc, err := LoadServer()
		So(err, ShouldBeNil)
		So(sc.Port, ShouldEqual, "9090")
	})
}

func TestServerAddr(t *testing.T) {
	Convey("LoadServer returns server subconfig", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Setenv(configPath, path)
		defer os.Unsetenv(configPath)
		sc, err := LoadServer()
		So(err, ShouldBeNil)
		So(sc.ServerAddr(), ShouldEqual, "127.0.0.1:9090")
	})
}

func TestDSN(t *testing.T) {
	Convey("LoadServer returns DSN", t, func() {
		path := writeTempConfig(t, sampleYAML)
		os.Setenv(configPath, path)
		defer os.Unsetenv(configPath)
		cfg, err := Load()
		So(err, ShouldBeNil)
		dsn := cfg.Storage.DSN()
		So(dsn, ShouldEqual, "postgresql://u:p@db:5432/n?sslmode=disable")
	})
}
