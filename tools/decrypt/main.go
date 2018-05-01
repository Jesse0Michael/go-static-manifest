package main

import (
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jesse0michael/go-static-manifest/builder"
)

func main() {
	keyFile := flag.String("key", "", "path to encryption key")
	iv := flag.String("iv", "", "hex string encryption iv")
	input := flag.String("input", "", "path to input file")
	output := flag.String("output", "", "path to output file")
	flag.Parse()

	if *keyFile == "" || *iv == "" || *input == "" || *output == "" {
		flag.Usage()
		os.Exit(1)
	}

	keyBytes, err := ioutil.ReadFile(*keyFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	encodedBytes := make([]byte, hex.EncodedLen(len(keyBytes)))
	_ = hex.Encode(encodedBytes, keyBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = builder.DecryptFile(*iv, string(encodedBytes), *input, *output)
	if err != nil {
		log.Fatal(err.Error())
	}
}
