package main

import (
	"fmt"
	"log"

	"github.com/agent-api/gsv"
)

type DeeperSchema struct {
	DeeperScore *gsv.IntSchema `json:"deeper_score"`
}

type MySchema struct {
	Name    *gsv.StringSchema `json:"name"`
	Address *gsv.StringSchema `json:"address"`
	Score   *gsv.IntSchema    `json:"score"`
	Deeper  *DeeperSchema     `json:"deeper"`
}

const jsonData = `{"name": "John", "address": "123 main st", "score": 99, "deeper": {"deeper_score": 1}}`

func main() {
	var schema MySchema
	schema.Name = gsv.String().Min(3).Max(50).Description("The name of the business")
	schema.Address = gsv.String().Description("The address of the business")

	schema.Deeper = &DeeperSchema{}
	schema.Deeper.DeeperScore = gsv.Int().Min(0).Max(99)

	_, err := gsv.Parse([]byte(jsonData), &schema)
	if err != nil {
		log.Fatal(err)
	}

	name, ok := schema.Name.Value()
	if !ok {
		log.Fatal("name is null")
	}

	address, ok := schema.Address.Value()
	if !ok {
		log.Fatal("address is null")
	}

	score, ok := schema.Score.Value()
	if !ok {
		log.Fatal("score is null")
	}

	deeperScore, ok := schema.Deeper.DeeperScore.Value()
	if !ok {
		log.Fatal("deeper score is null")
	}

	fmt.Println("Valid name:", name)
	fmt.Println("Valid address:", address)
	fmt.Println("Valid score:", score)
	fmt.Println("Valid deeper score:", deeperScore)
}
