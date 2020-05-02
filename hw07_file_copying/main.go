package main

import (
	"flag"
	"log"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

const ExitStatusError = 1

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Place your code here
	if from == "" || to == "" {
		flag.Usage()
		os.Exit(ExitStatusError)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatal(err)
	}
}
