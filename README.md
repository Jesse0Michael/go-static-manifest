# GO STATIC MANIFEST
Go Static Manifest is a small CLI tool, written in GO, that will download a HLS manifest to disk. Unlike other tools, this tool does not concatinate the segments into one file, but instead writes the M3U8's and segments to disk, re-writing the manifests to have the relative path, so the directory can be staically hosted. The tool also downloads and rewrites the relative paths for manifest encryption keys.

## Running

1. `go install github.com/jesse0michael/go-static-manifest`
2. `go-static-manifest {options} {manifest_url}`

```bash
Usage of ./bin/go-static-manifest:
  ./bin/go-static-manifest {options: optional} {manifest_url: required}
  -directory string
    	directory to write the manifest too (default "manifest")
```

# Additional Tools
Decrypt or Encrypt any file with a provided Key and IV.

1. `make build`
2. `go-(decrypt|encrypt)-file -key keyFile.key -iv 0x0000000000000000 -input inputFile -output outputFile`

```bash
Usage of go-(decrypt|encrypt)-file:
  -input string
    	path to input file
  -iv string
    	hex string encryption iv
  -key string
    	path to encryption key
  -output string
    	path to output file
```

