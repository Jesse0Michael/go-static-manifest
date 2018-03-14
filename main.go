package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jesse0michael/go-static-manifest/builder"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s {options: optional} {manifest_url: required}\n", os.Args[0])
		flag.PrintDefaults()
	}
	directory := flag.String("directory", "manifest", "directory to write the manifest too")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	manifest, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = builder.Build(manifest, *directory)
	if err != nil {
		log.Fatal(err.Error())
	}
}
