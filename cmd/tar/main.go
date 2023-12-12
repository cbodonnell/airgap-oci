package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func createTarGz(srcDir string, dest io.Writer) error {
	gw := gzip.NewWriter(dest)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	var files []string
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, path := range files {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// set modification time and other attributes to fixed values for deterministic output
		header.ModTime = time.Unix(0, 0)
		header.AccessTime = time.Unix(0, 0)
		header.ChangeTime = time.Unix(0, 0)
		header.Uid = 0
		header.Gid = 0
		header.Mode = int64(info.Mode().Perm())

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go source_directory destination_file.tar.gz")
		os.Exit(1)
	}

	srcDir := os.Args[1]
	destFile := os.Args[2]

	dest, err := os.Create(destFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer dest.Close()

	err = createTarGz(srcDir, dest)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s from %s\n", destFile, srcDir)
}
