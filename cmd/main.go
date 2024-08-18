package main

import (
	"log"

	"github.com/rizasghari/kalkan/cmd/kalkan"
)

func main() {
	if err := kalkan.New().Run(); err != nil {
		log.Fatalf("error running kalkan: %v", err)
	}
}
