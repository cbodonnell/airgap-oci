APP_SLUG ?= craig-helm

.PHONY: registry
registry:
	docker run --rm -p 5000:5000 --name registry registry:2

.PHONY: push
push:
	./scripts/push.sh

.PHONY: copy-postgres
copy-postgres:
	rm -rf layers/images/postgresql/16.1.0-debian-11-r15
	mkdir -p layers/images/postgresql/16.1.0-debian-11-r15
	skopeo copy --preserve-digests --multi-arch all docker://docker.io/bitnami/postgresql:16.1.0-debian-11-r15 dir:layers/images/postgresql/16.1.0-debian-11-r15

.PHONY: push-postgres
push-postgres:
	skopeo copy --dest-tls-verify=false --preserve-digests --multi-arch all dir:layers/images/postgresql/16.1.0-debian-11-r15 docker://localhost:5000/some-namespace/postgresql:16.1.0-debian-11-r15