package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var scriptTemplate = `#!/bin/bash

oras push \
  --artifact-type application/vnd.airgap.bundle \
  localhost:5001/my-team/airgap-bundle:v1 \
  bundle/airgap.yaml:application/vnd.airgap.metadata+yaml \
  bundle/app.tar.gz:application/vnd.airgap.application.tar+gzip \
  %s
`

func main() {
	blobsDIr := "bundle/images/docker/registry/v2/blobs"

	var filepaths []string

	err := filepath.WalkDir(blobsDIr, func(fp string, fi os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if !fi.IsDir() {
			filepaths = append(filepaths, fp)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf(scriptTemplate, strings.Join(filepaths, " \\\n  "))
}
