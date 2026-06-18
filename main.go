package main

import (
	"ballot-tool/internal/app"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// requires organization roles // outputs P- and O-memeber committees without voter
	tool := flag.String("tool", "ballots", "which tool to run")
	member := flag.Bool("member", false, "enable option")
	dev := flag.Bool("dev", false, "set true to use test-date from import tool")
	job := flag.String("job", "all", "choose what stats you want") //all, national, adoptions, norsok (all does not include norsok at the moment)
	nsOnly := flag.Bool("ns_only", false, "set to false to include all product types")
	//urn := flag.String("urn", "snv:proj:1973783", "urn of project to retrieve")
	from := flag.String("from", "2026-05-27", "publication date range begin")
	to := flag.String("to", "2026-06-03", "publication date range begin end")
	aktualitet := flag.Bool("aktualitet", false, "use flag to generate aktualitetsundersøkelse for current year")
	opt := flag.String("opt", "", "use to specify job variable")
	file := flag.String("file", "", "name and file extension of input file")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	switch *tool {
	case "ballots":
		if err := app.RunBallotTool(*member); err != nil {
			log.Fatalf("noe gikk galt: %s", err)
		}
	case "standards":
		if err := app.RunStandardsTool(*job, *from, *to, *file, *opt, *nsOnly, *aktualitet, *dev); err != nil {
			log.Fatalf("noe gikk galt: %s", err)
		}
	default:
		log.Fatalf("unknown tool: %s\n", *tool)
	}
}
