package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("./living.torrent")

	decodedValue, err := Decode(file)

	if err != nil {
		fmt.Println("error could not parse value", err)
	}

	jsonValue, err := json.MarshalIndent(decodedValue, " ", "  ")

	if err != nil {
		fmt.Println("error could not parse value", err)
	}

	fmt.Println(string(jsonValue))
}
