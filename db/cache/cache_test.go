package cache

import (
	"testing"

	"github.com/WhiCu/stgorders/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLRUCache_Basic(t *testing.T) {
	Convey("LRUCache Set and Get", t, func() {
		c, err := NewLRUCache[string, int](8, logger.NewNOPSlog())
		So(err, ShouldBeNil)
		So(c.Size(), ShouldEqual, 8)
		So(c.Set("a", 1), ShouldBeNil)
		v, err := c.Get("a")
		So(err, ShouldBeNil)
		So(v, ShouldEqual, 1)
	})

	Convey("LRUCache Get not found", t, func() {
		c, err := NewLRUCache[string, int](8, logger.NewNOPSlog())
		So(err, ShouldBeNil)
		_, err = c.Get("missing")
		So(err, ShouldEqual, ErrNotFound)
	})
}

func TestLRUCache_Map(t *testing.T) {
	Convey("LRUCache Map", t, func() {
		c, err := NewLRUCache[int, int](8, logger.NewNOPSlog())
		So(err, ShouldBeNil)
		for i := 0; i < 16; i++ {
			c.Set(i, i)
		}
		m := c.Map()
		So(len(m), ShouldEqual, 8)
	})
}

func TestErrSize(t *testing.T) {
	Convey("LRUCache ErrSize", t, func() {
		for i := 0; i < MinSize; i++ {
			_, err := NewLRUCache[int, int](i, logger.NewNOPSlog())
			So(err, ShouldNotBeNil)
		}
	})
}
