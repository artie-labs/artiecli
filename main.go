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

	command, err := internal.ParseCommand(os.Args)
	if err != nil {
		log.Fatalf("Failed to parse command: %v", err)
	}

	ctx := context.Background()
	artieClient := internal.NewArtieClient(os.Getenv("ARTIE_API_KEY"), os.Getenv("ARTIE_API_URL_OVERRIDE"))
	if err := command.Execute(ctx, artieClient); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
