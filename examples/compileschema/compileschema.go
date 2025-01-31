package main

import (
	"fmt"
	"log"

	"github.com/agent-api/gsv"
)

type UserSchema struct {
	Name    *gsv.StringSchema `json:"name"`
	Address *gsv.StringSchema `json:"address"`
}

func main() {
	var schema UserSchema
	schema.Name = gsv.String().Min(3).Max(50).Description("The name of the business")
	schema.Address = gsv.String().Description("The address of the business")

	compileSchemaOpts := &gsv.CompileSchemaOpts{
		SchemaTitle:       "User Schema",
		SchemaDescription: "This is my user schema",
	}
	s, err := gsv.CompileSchema(schema, compileSchemaOpts)
	if err != nil {
		log.Fatal("error during CompileSchema:", err)
	}

	fmt.Println("Compiled schema:", string(s))
}
