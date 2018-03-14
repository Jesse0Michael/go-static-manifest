# GO STATIC MANIFEST
Go Static Manifest is a small CLI tool, written in GO, that will download a HLS manifest to disk. Unlike other tools, this tool does not concatinate the segments into one file, but instead writes the M3U8's and segments to disk, re-writing the manifests to have the relative path, so the directory can be staically hosted.

## Running

1. `go install github.com/jesse0michael/go-static-manifest`
2. `go-static-manifest {options} {manifest_url}`

```bash
Usage of ./bin/go-static-manifest:
  ./bin/go-static-manifest {options: optional} {manifest_url: required}
  -directory string
    	directory to write the manifest too (default "manifest")
```

