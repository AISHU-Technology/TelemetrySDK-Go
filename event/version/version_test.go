package version

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVersion(t *testing.T) {
	convey.Convey("TestVersion", t, func() {
		convey.So(EventInstrumentationVersion, convey.ShouldEqual, "2.6.2")
	})
}
