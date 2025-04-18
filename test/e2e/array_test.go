package gsv_e2e_test

import (
	"github.com/agent-api/gsv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArraySchema", func() {
	Describe("Implements Schema interface", func() {
		// Compile time check for ArraySchema implementing the Schema interface
		schema := gsv.Array(gsv.String())
		var _ gsv.Schema = schema
	})

	Describe("New ArraySchema", func() {
		It("creates a new ArraySchema", func() {
			v := gsv.Array(gsv.String())
			Expect(v).ToNot(BeNil())
		})

		It("panics when element schema is nil", func() {
			Expect(func() {
				gsv.Array(nil)
			}).To(Panic())
		})
	})

	FDescribe("MinItems Validation", func() {
		DescribeTable("validates minimum items",
			func(value []string, min int, expectError bool) {
				v := gsv.Array(gsv.String()).MinItems(min)

				// Convert strings to interface{}
				items := make([]interface{}, len(value))
				for i, str := range value {
					items[i] = str
				}

				v.Set(items...)
				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MinItemsError))
					Expect(result.Errors[0].Expected).To(Equal(min))
					Expect(result.Errors[0].Actual).To(Equal(len(value)))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			Entry("valid length", []string{"hello", "world"}, 1, false),
			Entry("too few items", []string{}, 1, true),
			Entry("exact minimum", []string{"hello"}, 1, false),
		)
	})

	Describe("MaxItems Validation", func() {
		DescribeTable("validates maximum items",
			func(value []string, max int, expectError bool) {
				v := gsv.Array(gsv.String()).MaxItems(max)

				items := make([]interface{}, len(value))
				for i, str := range value {
					items[i] = str
				}

				err := v.Set(items)
				Expect(err).ToNot(HaveOccurred())

				result := v.Validate()

				if expectError {
					Expect(result.HasErrors()).To(BeTrue())
					Expect(result.Errors[0].Type).To(Equal(gsv.MaxItemsError))
					Expect(result.Errors[0].Expected).To(Equal(max))
					Expect(result.Errors[0].Actual).To(Equal(len(value)))
				} else {
					Expect(result.HasErrors()).To(BeFalse())
				}
			},
			Entry("valid length", []string{"hello"}, 2, false),
			Entry("too many items", []string{"hello", "world", "!"}, 2, true),
			Entry("exact maximum", []string{"hello", "world"}, 2, false),
		)
	})

	Describe("Element Validation", func() {
		It("validates each element using the element schema", func() {
			v := gsv.Array(gsv.String().Min(3))

			err := v.Set([]interface{}{"hi", "hello", "a"})
			Expect(err).ToNot(HaveOccurred())

			result := v.Validate()
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors).To(HaveLen(2)) // "hi" and "a" are too short
			Expect(result.Errors[0].Message).To(ContainSubstring("element 0"))
			Expect(result.Errors[1].Message).To(ContainSubstring("element 2"))
		})
	})

	Describe("Clone functionality", func() {
		It("creates an independent copy of the schema", func() {
			original := gsv.Array(gsv.String().Min(3)).MinItems(1).MaxItems(5)
			err := original.Set([]interface{}{"hello"})
			Expect(err).ToNot(HaveOccurred())

			cloned := original.Clone()

			// Modify original
			err = original.Set([]interface{}{"hi"})
			Expect(err).ToNot(HaveOccurred())

			// Validate cloned maintains its own state
			//val, ok := cloned.Value()
			//Expect(ok).To(BeTrue())
			//Expect(val).To(Equal([]interface{}{"hello"}))

			// Validate cloned maintains validation rules
			arrayClone := cloned.(*gsv.ArraySchema)
			err = arrayClone.Set([]interface{}{})
			Expect(err).ToNot(HaveOccurred())
			result := arrayClone.Validate()
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(gsv.MinItemsError))
		})
	})

	Describe("JSON Unmarshaling", func() {
		type TestNestedArraySchema struct {
			NestedTest *gsv.ArraySchema `json:"nested_test"`
		}

		type TestArraySchema struct {
			Test         *gsv.ArraySchema       `json:"test"`
			NestedSchema *TestNestedArraySchema `json:"nested_schema"`
		}

		It("correctly unmarshals nested array schemas", func() {
			jsonData := `{
                "test": ["hello", "world"],
                "nested_schema": {
                    "nested_test": ["nested", "array", "test"]
                }
            }`

			var schema TestArraySchema
			schema.Test = gsv.Array(gsv.String())
			schema.NestedSchema = &TestNestedArraySchema{
				NestedTest: gsv.Array(gsv.String()),
			}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).NotTo(HaveOccurred())

			val, ok := schema.Test.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(HaveLen(2))
			Expect(val[0]).To(Equal("hello"))

			val, ok = schema.NestedSchema.NestedTest.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(HaveLen(3))
			Expect(val[2]).To(Equal("test"))
		})

		It("handles optional arrays", func() {
			jsonData := `{"nested_schema": {"nested_test": ["test"]}}`

			var schema TestArraySchema
			schema.Test = gsv.Array(gsv.String()).Optional()
			schema.NestedSchema = &TestNestedArraySchema{
				NestedTest: gsv.Array(gsv.String()),
			}

			_, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).NotTo(HaveOccurred())

			val, ok := schema.Test.Value()
			Expect(ok).To(BeFalse())
			Expect(val).To(BeNil())
		})

		It("validates array elements during unmarshaling", func() {
			jsonData := `{"test": ["hi", "hello", "a"]}`

			var schema TestArraySchema
			schema.Test = gsv.Array(gsv.String().Min(3))

			result, err := gsv.Parse([]byte(jsonData), &schema)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.HasErrors()).To(BeTrue())

			// Should have validation errors for "hi" and "a"
			errors := result.HasErrors()
			Expect(errors).To(HaveLen(2))
		})
	})

	Describe("Value interface methods", func() {
		It("handles getValue and Value correctly", func() {
			schema := gsv.Array(gsv.String())

			// Test empty state
			_, ok := schema.Value()
			Expect(ok).To(BeFalse())

			// Set and get value
			testData := []interface{}{"hello", "world"}
			err := schema.Set(testData)
			Expect(err).ToNot(HaveOccurred())

			val, ok := schema.Value()
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(testData))
		})

		It("panics on invalid internal value type", func() {
			schema := gsv.Array(gsv.String())

			// This is a contrived example to test the panic - in real usage
			// this should never happen due to type checking in SetValue
			schemaAny := interface{}(schema)
			arraySchema := schemaAny.(*gsv.ArraySchema)

			// Deliberately set invalid internal value type
			err := arraySchema.Set("not an array")
			Expect(err).To(HaveOccurred())

			// Should panic when trying to get the value
			Expect(func() {
				arraySchema.Value()
			}).To(Panic())
		})
	})
})
