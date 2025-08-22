package app

import (
	"testing"

	"github.com/WhiCu/stgorders/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetLogger(t *testing.T) {
	Convey("getLogger builds logger without file", t, func() {
		cfg := &config.LoggerConfig{Level: "debug", Path: "", Size: 1, Compress: false}
		l := getLogger(cfg)
		So(l, ShouldNotBeNil)
	})
}
