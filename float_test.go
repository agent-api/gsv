package gsv_test

import (
	"github.com/agent-api/gsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Float32Schema", func() {
	It("Implements the Schema interface", func() {
		schema := gsv.Float32()
		var _ gsv.Schema = schema
	})

	Describe("New Float32Schema", func() {
		It("creates a new Float32Schema", func() {
			v := gsv.Float32()
			Expect(v).ToNot(BeNil())
		})
	})

	Describe("Min Validation", func() {
		DescribeTable("validates minimum float32 value",
			func(value float32, min float32, message string, expectError bool) {
				v := gsv.Float32().Min(min)
				v.Set(value)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MinNumberError))
					Expect(result.Errors[0].Message).To(Equal(message))
					Expect(result.Errors[0].Expected).To(Equal(min))
					Expect(result.Errors[0].Actual).To(Equal(value))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			// Happy paths
			Entry("greater than", float32(5.123), float32(3.123), "", false),
			Entry("exact size", float32(3.123), float32(3.123), "", false),

			// Error cases
			Entry("smaller than", float32(1.123), float32(3.123), "must be at least 3.123", true),
		)

		It("supports custom error messages", func() {
			customMsg := "too small!"
			v := gsv.Float32().Min(float32(3.123), gsv.ValidationOptions{Message: customMsg})
			v.Set(float32(1.123))
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MinNumberError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	Describe("Min Validation", func() {
		DescribeTable("validates maximum float32 value",
			func(value float32, max float32, message string, expectError bool) {
				v := gsv.Float32().Max(max)
				v.Set(value)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MaxNumberError))
					Expect(result.Errors[0].Message).To(ContainSubstring(message))
					Expect(result.Errors[0].Expected).To(Equal(max))
					Expect(result.Errors[0].Actual).To(Equal(value))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			// Happy paths
			Entry("less than", float32(3.123), float32(5.123), "", false),
			Entry("exact size", float32(3.123), float32(3.123), "", false),

			// Error cases
			Entry("greater than", float32(5.123), float32(3.123), "must not exceed: 3.123", true),
		)

		It("supports custom error messages", func() {
			customMsg := "too big!"
			v := gsv.Float32().Max(float32(3.123), gsv.ValidationOptions{Message: customMsg})
			v.Set(float32(5.123))
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MaxNumberError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	XDescribe("Min and Max Combined", func() {
		DescribeTable("validates both min and max",
			func(value, min, max float32, expectedErrors map[string]string) {
				v := gsv.Float32().Min(min).Max(max)
				v.Set(value)
				result := v.Validate()

				if len(expectedErrors) > 0 {
					Expect(result.HasErrors()).To(BeTrue())
					for _, err := range result.Errors {
						if msg, exists := expectedErrors[string(err.Type)]; exists {
							Expect(err.Message).To(Equal(msg))
						}
					}
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			Entry("valid length", float32(4.123), float32(3.123), float32(5.123), nil),
			Entry("too small", float32(1.123), float32(3.123), float32(5.123), map[string]string{
				string(gsv.MinNumberError): "must be at least 3.123",
			}),
			Entry("too big", float32(6.123), float32(3.123), float32(5.123), map[string]string{
				string(gsv.MinNumberError): "must be at no more than 5.123",
			}),
			Entry("exact minimum", float32(3.123), float32(3.123), float32(5.123), nil),
			Entry("exact maximum", float32(5.123), float32(3.123), float32(5.123), nil),
		)
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

		XIt("handles null schema pofloat32ers when they're not instantiated", func() {
			jsonData := `{"nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Float64Schema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Float64()
		var _ gsv.Schema = schema
	})
})
