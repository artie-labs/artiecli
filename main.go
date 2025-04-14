package main

import (
	"log"

	"github.com/artie-labs/transfer/lib/environ"
)

func main() {
	if err := environ.MustGetEnv("ARTIE_API_KEY"); err != nil {
		log.Fatalf("ARTIE_API_KEY is not set: %v", err)
	}
}
