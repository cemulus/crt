package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cemulus/crt/repository"
	"github.com/cemulus/crt/result"
)

var (
	filename  = flag.String("o", "", "")
	subdomain = flag.Bool("s", false, "")
	expired   = flag.Bool("e", false, "")
	limit     = flag.Int("l", 1000, "")
	jsonOut   = flag.Bool("json", false, "")
	csvOut    = flag.Bool("csv", false, "")
)

var usage = `Usage: crt [options...] <domain name>

Options:
  -o <path> Output file path. Write to file instead of stdout.
  -s        Enumerate subdomains.
  -e        Exclude expired certificates.
  -l <int>  Limit the number of results. (default: 1000) 
  -json     Turn results to JSON.
  -csv      Turn results to CSV.

Examples:
  crt example.com
  crt -s -e example.com
  crt -o logs.json -json example.com
  crt -csv -o logs.csv -l 15 example.com
`

func Execute() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	domain := flag.Args()[0]
	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	repo, err := repository.New()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	var res result.Printer

	if *subdomain {
		res, err = repo.GetSubdomains(domain, *expired, *limit)
	} else {
		res, err = repo.GetCertLogs(domain, *expired, *limit)
	}
	if err != nil {
		log.Fatal(err)
	}

	if res.Size() == 0 {
		fmt.Println("Found no results.")
		os.Exit(0)
	}

	var out []byte

	if *jsonOut {
		out, err = res.JSON()
	} else if *csvOut {
		out, err = res.CSV()
	} else {
		out = res.Table()
	}

	if err != nil {
		log.Fatal(err)
	}

	if *filename == "" {
		fmt.Print(string(out))
		os.Exit(0)
	}

	file, err := os.Create(*filename)
	if err != nil {
		log.Fatal("failed to create output file:", err)
	}
	defer file.Close()

	if _, err = file.Write(out); err != nil {
		log.Fatal("failed to write to file:", err)
	}
}
