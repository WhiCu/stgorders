package logger

import (
	"context"
	"log/slog"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMustInitLogger(t *testing.T) {
	Convey("MustInitLogger returns handler for known levels", t, func() {
		for _, lvl := range []string{"debug", "info", "warn", "error"} {
			h := MustInitLogger(lvl)
			So(h, ShouldNotBeNil)
		}
	})

	Convey("MustInitLogger panics on invalid level", t, func() {
		So(func() { MustInitLogger("invalid") }, ShouldPanic)
	})
}

func TestNOPHandler(t *testing.T) {
	Convey("NOPHandler disables all logs and returns itself on With*", t, func() {
		h := NewNOPHandler()
		So(h.Enabled(context.TODO(), slog.LevelDebug), ShouldBeFalse)
		So(h.Handle(context.TODO(), slog.Record{}), ShouldBeNil)
		So(h.WithAttrs([]slog.Attr{}), ShouldResemble, h)
		So(h.WithGroup(""), ShouldResemble, h)
		l := NewNOPSlog()
		So(l, ShouldNotBeNil)
	})
}
