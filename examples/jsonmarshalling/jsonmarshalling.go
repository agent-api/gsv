package main

import (
	"encoding/json"
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
	schema.Name = gsv.String().Set("John")
	schema.Address = gsv.String().Set("123 Main St")

	s, err := json.Marshal(schema)
	if err != nil {
		log.Fatal("error during marshalling:", err)
	}

	fmt.Println("Schema marshalled:", string(s))
}
