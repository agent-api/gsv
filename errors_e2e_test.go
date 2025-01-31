package gsv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/agent-api/gsv"
)

var _ = Describe("Schema Errors", func() {
	Context("Validation Error Details", func() {
		It("provides detailed error information", func() {
			schema := gsv.String().Min(5).Max(10).Set("abc")

			result := schema.Validate()
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors).To(HaveLen(1))

			err := result.Errors[0]
			Expect(err.Type).To(Equal(gsv.MinStringLengthError))
			Expect(err.Expected).To(Equal(5))
			Expect(err.Actual).To(Equal(3))
		})
	})
})
