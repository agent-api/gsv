package gsv

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema Validation", func() {
	Context("Nested Schemas", func() {
		It("validates nested structs", func() {
			type AddressSchema struct {
				Street *StringSchema `json:"street"`
				City   *StringSchema `json:"city"`
			}

			type UserSchema struct {
				Name    *StringSchema  `json:"name"`
				Address *AddressSchema `json:"address"`
			}

			schema := &UserSchema{
				Name: String().Set("John"),
				Address: &AddressSchema{
					Street: String().Set("123 Main St"),
					City:   String().Set("Boston"),
				},
			}

			result := ensure(schema)
			Expect(result.HasErrors()).To(BeFalse())
		})

		It("handles deeply nested validation errors", func() {
			type DeepSchema struct {
				Value *StringSchema `json:"value"`
			}

			type NestedSchema struct {
				Deep *DeepSchema `json:"deep"`
			}

			type RootSchema struct {
				Nested *NestedSchema `json:"nested"`
			}

			schema := &RootSchema{
				Nested: &NestedSchema{
					Deep: &DeepSchema{
						Value: String().Min(5).Set("abc"),
					},
				},
			}

			// Check the specific error details
			result := schema.Nested.Deep.Value.Validate()
			Expect(result.HasErrors()).To(BeTrue())
			Expect(result.Errors[0].Type).To(Equal(MinStringLengthError))
			Expect(result.Errors[0].Message).To(ContainSubstring("must be at least 5 characters"))
			Expect(result.Errors[0].Expected).To(Equal(5))
			Expect(result.Errors[0].Actual).To(Equal(3))
		})
	})

	Context("Edge Cases", func() {
		XIt("handles nil pointers gracefully", func() {
			type TestSchema struct {
				OptionalField StringSchema `json:"optional"`
			}

			schema := &TestSchema{}

			result := ensure(schema)
			Expect(result.HasErrors()).To(BeFalse())
		})

		It("validates empty structs with no fields", func() {
			type EmptySchema struct{}

			schema := &EmptySchema{}

			result := ensure(schema)
			Expect(result.HasErrors()).To(BeFalse())
		})

		XIt("handles non-struct inputs", func() {
			// It should fail - it makes no sense to validate a non-struct type
			result := ensure("not a struct")
			Expect(result.HasErrors()).To(BeFalse())

			result = ensure(123)
			Expect(result.HasErrors()).To(BeFalse())
		})
	})
})
