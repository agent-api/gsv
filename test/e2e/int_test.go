package gsv_e2e_test

import (
	"github.com/agent-api/gsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntSchema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Int()
		var _ gsv.Schema = schema
	})

	Describe("New IntSchema", func() {
		It("creates a new IntSchema", func() {
			v := gsv.Int()
			Expect(v).ToNot(BeNil())
		})
	})

	Describe("Min Validation", func() {
		DescribeTable("validates minimum int value",
			func(value int, min int, message string, expectError bool) {
				v := gsv.Int().Min(min)
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
			Entry("greater than", 5, 3, "", false),
			Entry("exact size", 3, 3, "", false),

			// Error cases
			Entry("smaller than", 1, 3, "must be at least 3", true),
		)

		It("supports custom error messages", func() {
			customMsg := "too small!"
			v := gsv.Int().Min(3, gsv.ValidationOptions{Message: customMsg})
			v.Set(1)
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MinNumberError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	Describe("Min Validation", func() {
		DescribeTable("validates maximum int value",
			func(value int, max int, message string, expectError bool) {
				v := gsv.Int().Max(max)
				v.Set(value)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MaxNumberError))
					Expect(result.Errors[0].Message).To(Equal(message))
					Expect(result.Errors[0].Expected).To(Equal(max))
					Expect(result.Errors[0].Actual).To(Equal(value))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			// Happy paths
			Entry("less than", 3, 5, "", false),
			Entry("exact size", 3, 3, "", false),

			// Error cases
			Entry("greater than", 5, 3, "must not exceed: 3", true),
		)

		It("supports custom error messages", func() {
			customMsg := "too big!"
			v := gsv.Int().Max(3, gsv.ValidationOptions{Message: customMsg})
			v.Set(5)
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MaxNumberError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	XDescribe("Min and Max Combined", func() {
		DescribeTable("validates both min and max constraints",
			func(value string, min, max int, expectedErrors map[string]string) {
				v := gsv.String().Min(min).Max(max)
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
			Entry("valid length", "hello", 3, 10, nil),
			Entry("too short", "hi", 3, 10, map[string]string{
				string(gsv.MinNumberError): "must be at least 3 characters long",
			}),
			Entry("too long", "hello world", 3, 5, map[string]string{
				string(gsv.MaxNumberError): "must be at most 5 characters long",
			}),
			Entry("exact minimum", "hey", 3, 5, nil),
			Entry("exact maximum", "hello", 3, 5, nil),
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

		XIt("handles null schema pointers when they're not instantiated", func() {
			jsonData := `{"nested_schema": {"nested_test": "nested test"}}`

			var schema TestStringSchema
			schema.NestedSchema = &TestNestedStringSchema{}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Int8Schmea", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Int8()
		var _ gsv.Schema = schema
	})
})

var _ = Describe("Int16Schema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Int16()
		var _ gsv.Schema = schema
	})
})

var _ = Describe("Int32Schema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Int32()
		var _ gsv.Schema = schema
	})
})

var _ = Describe("Int64Schema", func() {
	It("implements the Schema interface", func() {
		schema := gsv.Int64()
		var _ gsv.Schema = schema
	})
})
