.PHONY: bundle-registry
bundle-registry:
	docker run --rm -p 5000:5000 -v $(PWD)/bundle/images:/var/lib/registry --name bundle-registry registry:2

.PHONY: storage-registry
storage-registry:
	docker run --rm -p 5001:5000 --name storage-registry registry:2

.PHONY: onprem-registry
onprem-registry:
	docker run --rm -p 5002:5000 --name onprem-registry registry:2

.PHONY: copy-images
copy-images:
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch all docker://docker.io/library/postgres:15.5-alpine docker://localhost:5000/postgres:15.5-alpine
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch all docker://docker.io/library/postgres:16.1-alpine docker://localhost:5000/postgres:16.1-alpine

.PHONY: script
script:
	go run ./cmd/main.go > ./scripts/artifact.sh
	chmod +x ./scripts/artifact.sh

.PHONY: artifact
artifact:
	./scripts/artifact.sh
