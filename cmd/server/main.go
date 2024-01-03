package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

var progress = make(map[string]int64)
var checksums = make(map[string]string)

func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	// client-side idenfifier for this download
	downloadUUID := r.URL.Query().Get("uuid")
	progress[downloadUUID] = 0
	checksums[downloadUUID] = "pending..."

	src, err := remote.NewRepository("localhost:5001/my-company/airgap-bundle")
	if err != nil {
		fmt.Printf("Error creating remote repository client: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	src.PlainHTTP = true // for insecure registry

	reference := "v1"
	// reference := "sha256:3c6f36667255d9f0798bd0d082e69a7778a1fa3f11d0366afa67f52949427132"
	descriptor, err := src.Resolve(r.Context(), reference)
	if err != nil {
		fmt.Printf("Error resolving artifact descriptor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := content.FetchAll(r.Context(), src, descriptor)
	if err != nil {
		fmt.Printf("Error fetching artifact content: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var manifest ocispec.Manifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		fmt.Printf("Error unmarshalling manifest: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=my-application-v1.airgap")
	w.Header().Set("Content-Type", "application/gzip")
	// we don't know the total size of the download as it's compressed on the fly
	// this will result in an indeterminate progress bar in the browser indicating
	// that something is downloading, but with no progress indication
	w.Header().Set("Transfer-Encoding", "chunked")

	hasher := sha256.New()
	// we can't know the checksum until the entire file has been written
	// and the gzip writer has been closed
	defer func() {
		checksum := hasher.Sum(nil)
		checksumHex := hex.EncodeToString(checksum)
		checksums[downloadUUID] = checksumHex
	}()

	multiWriter := io.MultiWriter(w, hasher)

	gzipWriter := gzip.NewWriter(multiWriter)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	totalSize := int64(0)
	for _, layer := range manifest.Layers {
		totalSize += layer.Size
	}
	processedSize := int64(0)

	for _, layer := range manifest.Layers {
		err := copyLayerToTarWriter(r.Context(), src.Blobs(), layer, tarWriter)
		if err != nil {
			fmt.Printf("Error copying layer to tar writer: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		processedSize += layer.Size
		progress[downloadUUID] = processedSize * 100 / totalSize
	}
}

func copyLayerToTarWriter(ctx context.Context, blobStore registry.BlobStore, layer ocispec.Descriptor, dst *tar.Writer) error {
	rc, err := blobStore.Fetch(ctx, layer)
	if err != nil {
		return fmt.Errorf("error fetching layer content: %w", err)
	}
	defer rc.Close()

	header := &tar.Header{
		Name: layer.Annotations[ocispec.AnnotationTitle],
		Mode: 0600,
		Size: layer.Size,
	}

	err = dst.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("error writing tar header: %w", err)
	}

	_, err = io.Copy(dst, rc)
	if err != nil {
		return fmt.Errorf("error copying layer content: %w", err)
	}

	return nil
}

func handleProgressRequest(w http.ResponseWriter, r *http.Request) {
	downloadUUID := r.URL.Query().Get("uuid")
	progress, ok := progress[downloadUUID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d", progress)))
}

func handleChecksumRequest(w http.ResponseWriter, r *http.Request) {
	downloadUUID := r.URL.Query().Get("uuid")
	checksum, ok := checksums[downloadUUID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(checksum))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/index.html")
	})

	http.HandleFunc("/download", handleDownloadRequest)
	http.HandleFunc("/progress", handleProgressRequest)
	http.HandleFunc("/checksum", handleChecksumRequest)

	port := 8080
	fmt.Printf("Server started on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
