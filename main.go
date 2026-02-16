package main

import (
	"ballot-tool/internal/app"
	"flag"
	"log"
)

func main() {
	// requires organization roles // outputs P- and O-memeber committees without voter
	member := flag.Bool("opt", false, "enable option")

	if err := app.Run(*member); err != nil {
		log.Fatalf("noe gikk galt %s", err)
	}
}
