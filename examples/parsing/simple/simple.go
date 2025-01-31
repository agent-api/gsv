package main

import (
	"fmt"
	"log"

	"github.com/agent-api/gsv"
)

type MySchema struct {
	Field *gsv.StringSchema `json:"field"`
}

const jsonData = `{"field": "example field data"}`

func main() {
	var schema MySchema
	schema.Field = gsv.String()

	_, err := gsv.Parse([]byte(jsonData), &schema)
	if err != nil {
		log.Fatal(err)
	}

	name, ok := schema.Field.Value()
	if !ok {
		log.Fatal("field is null")
	}

	fmt.Println("Value for field:", name)
}
