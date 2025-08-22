package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNOPCache(t *testing.T) {
	Convey("NOPCache basic behavior", t, func() {
		c := NewNOPCache[string](0)
		So(c.Size(), ShouldEqual, 0)
		_, err := c.Get("x")
		So(err, ShouldEqual, ErrNotFound)
		So(c.Set("x", 1), ShouldEqual, ErrSetCache)
	})
}
