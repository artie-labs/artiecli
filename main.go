package main

import (
	"context"
	"log"
	"os"

	"github.com/artie-labs/artiecli/internal"
	"github.com/artie-labs/transfer/lib/environ"
)

func main() {
	if err := environ.MustGetEnv("ARTIE_API_KEY"); err != nil {
		log.Fatalf("ARTIE_API_KEY is not set: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("no command provided")
	}

	ctx := context.Background()
	artieClient := internal.NewArtieClient(os.Getenv("ARTIE_API_KEY"), os.Getenv("ARTIE_API_URL_OVERRIDE"))

	switch os.Args[1] {
	case "list-deployments":
		if err := artieClient.ListDeployments(ctx); err != nil {
			log.Fatalf("failed to list deployments: %v", err)
		}
	}
}
