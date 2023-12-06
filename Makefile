APP_SLUG ?= craig-helm

.PHONY: registry
registry:
	docker run --rm -p 5000:5000 -v $(PWD)/layers/images:/var/lib/registry --name registry registry:2

.PHONY: push
push:
	./scripts/push.sh

.PHONY: copy-images
copy-images:
	rm -rf layers/images
	mkdir -p layers/images
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch system docker://docker.io/library/postgres:15.5 docker://localhost:5000/postgres:15.5
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch system docker://docker.io/library/postgres:16.1 docker://localhost:5000/postgres:16.1

.PHONY: push-postgres
push-postgres:
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch system dir:layers/images/postgresql/16.1.0-debian-11-r15 docker://localhost:5000/some-namespace/postgresql:16.1.0-debian-11-r15