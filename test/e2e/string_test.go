package gsv_e2e_test

import (
	"github.com/agent-api/gsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringSchema", func() {
	Describe("Implements Schema intercace", func() {
		// Compile time check for StringSchema implementing the Schema interface
		schema := gsv.String()
		var _ gsv.Schema = schema

	})

	Describe("New StringSchema", func() {
		It("creates a new StringSchema", func() {
			v := gsv.String()
			Expect(v).ToNot(BeNil())
		})
	})

	Describe("Min Validation", func() {
		DescribeTable("validates minimum length",
			func(value string, min int, message string, expectError bool) {
				v := gsv.String().Min(min)
				v.Set(value)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MinStringLengthError))
					if message != "" {
						Expect(result.Errors[0].Message).To(Equal(message))
					}
					Expect(result.Errors[0].Expected).To(Equal(min))
					Expect(result.Errors[0].Actual).To(Equal(len(value)))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			Entry("valid length", "hello", 3, "", false),
			Entry("too short", "hi", 3, "must be at least 3 characters long", true),
			Entry("empty string", "", 1, "must be at least 1 characters long", true),
			Entry("exact minimum", "hey", 3, "", false),
		)

		It("supports custom error messages", func() {
			customMsg := "too short!"
			v := gsv.String().Min(3, gsv.ValidationOptions{Message: customMsg})
			v.Set("hi")
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MinStringLengthError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	Describe("Max Validation", func() {
		DescribeTable("validates maximum length",
			func(value string, max int, message string, expectError bool) {
				v := gsv.String().Max(max)
				v.Set(value)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MaxStringLengthError))
					if message != "" {
						Expect(result.Errors[0].Message).To(Equal(message))
					}
					Expect(result.Errors[0].Expected).To(Equal(max))
					Expect(result.Errors[0].Actual).To(Equal(len(value)))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			Entry("valid length", "hello", 10, "", false),
			Entry("too long", "hello world", 5, "must be at most 5 characters long", true),
			Entry("exact maximum", "hello", 5, "", false),
		)

		It("supports custom error messages", func() {
			customMsg := "too long!"
			v := gsv.String().Max(5, gsv.ValidationOptions{Message: customMsg})
			v.Set("hello world")
			result := v.Validate()

			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MaxStringLengthError))
			Expect(result.Errors[0].Message).To(Equal(customMsg))
		})
	})

	Describe("Min and Max Combined", func() {
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
				string(gsv.MinStringLengthError): "must be at least 3 characters long",
			}),
			Entry("too long", "hello world", 3, 5, map[string]string{
				string(gsv.MaxStringLengthError): "must be at most 5 characters long",
			}),
			Entry("exact minimum", "hey", 3, 5, nil),
			Entry("exact maximum", "hello", 3, 5, nil),
		)
	})

	Describe("JSON Unmarshaling", func() {
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

			result, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.HasErrors()).To(BeTrue())
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
