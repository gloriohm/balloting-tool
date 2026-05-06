package main

import (
	"ballot-tool/internal/app"
	"flag"
	"log"
)

func main() {
	// requires organization roles // outputs P- and O-memeber committees without voter
	member := flag.Bool("opt", false, "enable option")
	tool := flag.String("tool", "ballots", "which tool to run")
	job := flag.String("job", "all", "choose what stats you want") //all, national, adoptions, norsok (all does not include norsok at the moment)
	nsOnly := flag.Bool("ns_only", false, "set to false to include all product types")

	flag.Parse()

	switch *tool {
	case "ballots":
		if err := app.RunBallotTool(*member); err != nil {
			log.Fatalf("noe gikk galt: %s", err)
		}
	case "standards":
		if err := app.RunStandardsTool(*job, *nsOnly); err != nil {
			log.Fatalf("noe gikk galt: %s", err)
		}
	/* case "table":
	if err := app.RunTableTool(); err != nil {
		log.Fatalf("noe gikk galt: %s", err)
	} */
	default:
		log.Fatalf("unknown tool: %s\n", *tool)
	}
}
