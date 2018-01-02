// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package udp

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestAlwaysValid(t *testing.T) {
	validator := AlwaysValid()

	a := assertions.New(t)
	a.So(validator.ValidUplink(Packet{}), should.BeTrue)
	a.So(validator.ValidDownlink(Packet{}), should.BeTrue)
}
