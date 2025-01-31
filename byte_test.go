package gsv_test

import (
	"github.com/agent-api/gsv"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("ByteSchema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Byte()
		var _ gsv.Schema = schema
	})
})
