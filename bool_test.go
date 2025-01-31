package gsv_test

import (
	"github.com/agent-api/gsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BoolSchema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Bool()
		var _ gsv.Schema = schema
	})

	Describe("BoolSchema())", func() {
		It("creates a new boolSchema", func() {
			v := gsv.Int()
			Expect(v).ToNot(BeNil())
		})
	})

	XDescribe("JSON Unmarshaling", func() {
		type TestNestedStringSchema struct {
			NestedTest *gsv.StringSchema `json:"nested_test"`
		}

		type TestStringSchema struct {
			Test         *gsv.StringSchema       `json:"test"`
			NestedSchema *TestNestedStringSchema `json:"nested_schema"`
		}

		It("correctly unmarshals nested string schemas", func() {
			jsonData := `{"test": "testing json data", "nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).NotTo(HaveOccurred())

			val, ok := schema.Test.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("testing json data"))

			val, ok = schema.NestedSchema.NestedTest.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("nested test"))
		})

		It("handles optional fields correctly", func() {
			jsonData := `{"nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.Test = gsv.String().Optional()
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).NotTo(HaveOccurred())

			val, ok := schema.Test.Value()
			Expect(ok).To(BeFalse())
			Expect(val).To(Equal(""))

			val, ok = schema.NestedSchema.NestedTest.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("nested test"))
		})

		It("enforces required fields", func() {
			jsonData := `{"nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.Test = gsv.String()
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).To(HaveOccurred())
		})

		XIt("handles null schema pointers when they're not instantiated", func() {
			jsonData := `{"nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).To(HaveOccurred())
		})
	})
})
