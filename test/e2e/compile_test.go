package gsv_e2e_test

import (
	"encoding/json"

	"github.com/agent-api/gsv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema Compiler", func() {
	Context("when compiling string schemas", func() {
		type BasicStringSchema struct {
			Name *gsv.StringSchema `json:"name"`
		}

		type ComplexStringSchema struct {
			Required     *gsv.StringSchema `json:"required"`
			Optional     *gsv.StringSchema `json:"optional"`
			WithLengths  *gsv.StringSchema `json:"withLengths"`
			Description  *gsv.StringSchema `json:"description"`
			NoJsonTag    *gsv.StringSchema
			IgnoredField *gsv.StringSchema `json:"-"`
		}

		It("should compile a basic string schema", func() {
			schema := BasicStringSchema{
				Name: gsv.String().Description("The user's name"),
			}

			result, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"basic", "A basic schema"})
			Expect(err).NotTo(HaveOccurred())

			var jsonSchema map[string]interface{}
			err = json.Unmarshal(result, &jsonSchema)
			Expect(err).NotTo(HaveOccurred())

			// Verify schema structure
			Expect(jsonSchema["title"]).To(Equal("basic"))
			Expect(jsonSchema["description"]).To(Equal("A basic schema"))
			Expect(jsonSchema["type"]).To(Equal("object"))

			properties := jsonSchema["properties"].(map[string]interface{})
			Expect(properties).To(HaveKey("name"))

			nameProperty := properties["name"].(map[string]interface{})
			Expect(nameProperty["type"]).To(Equal("string"))
			Expect(nameProperty["description"]).To(Equal("The user's name"))
		})

		It("should compile a nested schema", func() {
			type NestedStringSchema struct {
				Nested *BasicStringSchema `json:"nested"`
			}

			schema := &NestedStringSchema{
				Nested: &BasicStringSchema{
					Name: gsv.String().Description("The user's name"),
				},
			}

			result, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"basic_nested", "A basic nested schema"})
			Expect(err).NotTo(HaveOccurred())

			var jsonSchema map[string]interface{}
			err = json.Unmarshal(result, &jsonSchema)
			Expect(err).NotTo(HaveOccurred())

			// Verify schema structure
			Expect(jsonSchema["title"]).To(Equal("basic_nested"))
			Expect(jsonSchema["description"]).To(Equal("A basic nested schema"))
			Expect(jsonSchema["type"]).To(Equal("object"))

			properties := jsonSchema["properties"].(map[string]interface{})
			Expect(properties).To(HaveKey("nested"))

			// Verify nested object structure
			nestedProperty := properties["nested"].(map[string]interface{})
			Expect(nestedProperty["type"]).To(Equal("object"))

			// Verify nested properties
			nestedProperties := nestedProperty["properties"].(map[string]interface{})
			Expect(nestedProperties).To(HaveKey("name"))

			// Verify nested name field
			nameProperty := nestedProperties["name"].(map[string]interface{})
			Expect(nameProperty["type"]).To(Equal("string"))
			Expect(nameProperty["description"]).To(Equal("The user's name"))

			// Verify required fields
			required := jsonSchema["required"].([]interface{})
			Expect(required).To(ContainElement("nested"))

			// Verify nested required fields
			nestedRequired := nestedProperty["required"].([]interface{})
			Expect(nestedRequired).To(ContainElement("name"))
		})

		It("should compile a complex string schema", func() {
			schema := ComplexStringSchema{
				Required:     gsv.String(),
				Optional:     gsv.String().Optional(),
				WithLengths:  gsv.String().Min(5).Max(10),
				Description:  gsv.String().Description("A field with description"),
				NoJsonTag:    gsv.String(),
				IgnoredField: gsv.String(),
			}

			result, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"complex", "A complex schema"})
			Expect(err).NotTo(HaveOccurred())

			var jsonSchema map[string]interface{}
			err = json.Unmarshal(result, &jsonSchema)
			Expect(err).NotTo(HaveOccurred())

			properties := jsonSchema["properties"].(map[string]interface{})

			// Check required field
			Expect(properties).To(HaveKey("required"))
			requiredProperty := properties["required"].(map[string]interface{})
			Expect(requiredProperty["type"]).To(Equal("string"))

			// Check optional field
			Expect(properties).To(HaveKey("optional"))
			optionalProperty := properties["optional"].(map[string]interface{})
			Expect(optionalProperty["type"]).To(Equal("string"))

			// Check field with lengths
			Expect(properties).To(HaveKey("withLengths"))
			lengthsProperty := properties["withLengths"].(map[string]interface{})
			Expect(lengthsProperty["minLength"]).To(Equal(float64(5)))
			Expect(lengthsProperty["maxLength"]).To(Equal(float64(10)))

			// Check description field
			Expect(properties).To(HaveKey("description"))
			descProperty := properties["description"].(map[string]interface{})
			Expect(descProperty["description"]).To(Equal("A field with description"))

			// Check that fields without json tags or with "-" are not included
			Expect(properties).NotTo(HaveKey("NoJsonTag"))
			Expect(properties).NotTo(HaveKey("IgnoredField"))

			// Check required fields array
			required := jsonSchema["required"].([]interface{})
			Expect(required).To(ContainElement("required"))
			Expect(required).To(ContainElement("withLengths"))
			Expect(required).To(ContainElement("description"))
			Expect(required).NotTo(ContainElement("optional"))
		})

		It("should error on nil schema fields", func() {
			schema := BasicStringSchema{
				Name: nil,
			}

			_, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"nil_fields", "Schema with nil field"})
			Expect(err).To(HaveOccurred())
		})

		Context("when compiling invalid schemas", func() {
			type InvalidSchema struct {
				NotASchema string `json:"invalid"`
			}

			It("should return an error for unsupported field types", func() {
				// Arrange
				schema := InvalidSchema{
					NotASchema: "this is not a schema",
				}

				// Act
				result, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"invalid", "Invalid schema"})

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported schema type"))
				Expect(result).To(BeNil())
			})
		})
	})

	Context("when handling edge cases", func() {
		type EdgeCaseSchema struct {
			EmptyDescription *gsv.StringSchema `json:"emptyDesc"`
			ZeroLength       *gsv.StringSchema `json:"zeroLength"`
		}

		It("should handle edge cases correctly", func() {
			// Arrange
			schema := EdgeCaseSchema{
				EmptyDescription: gsv.String().Description(""),
				ZeroLength:       gsv.String().Min(0).Max(0),
			}

			// Act
			result, err := gsv.CompileSchema(schema, &gsv.CompileSchemaOpts{"edge_cases", "Edge case schema"})

			// Assert
			Expect(err).NotTo(HaveOccurred())

			var jsonSchema map[string]interface{}
			err = json.Unmarshal(result, &jsonSchema)
			Expect(err).NotTo(HaveOccurred())

			properties := jsonSchema["properties"].(map[string]interface{})

			// Check empty description handling
			emptyDescProperty := properties["emptyDesc"].(map[string]interface{})
			Expect(emptyDescProperty).NotTo(HaveKey("description"))

			// Check zero length handling
			zeroLengthProperty := properties["zeroLength"].(map[string]interface{})
			Expect(zeroLengthProperty["minLength"]).To(Equal(float64(0)))
			Expect(zeroLengthProperty["maxLength"]).To(Equal(float64(0)))
		})
	})
})
