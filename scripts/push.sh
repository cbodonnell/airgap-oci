#!/bin/bash

# NOTES
# - `vnd.airgap.docker.` is used in place of `vnd.docker.` for image manifests because they would not be uploaded as blobs otherwise

oras push \
  --artifact-type application/vnd.airgap.bundle \
  --config metadata/airgap.json:application/vnd.airgap.metadata.v1+json \
  registry.replicated.com/${APP_SLUG}/airgap-artifact:v1 \
  layers/config/postgresql-13.2.23.tgz:application/vnd.cncf.helm.chart.content.v1.tar+gzip \
  layers/config/postgresql.yaml:application/vnd.airgap.config.content.v1+yaml \
  layers/config/config.yaml:application/vnd.airgap.config.content.v1+yaml \
  layers/config/replicated-app.yaml:application/vnd.airgap.config.content.v1+yaml \
  layers/images/postgresql/16.1.0-debian-11-r15/manifest.json:application/vnd.airgap.docker.distribution.manifest.list.v2+json \
  layers/images/postgresql/16.1.0-debian-11-r15/9d1334877f8ca27be121f122dea7e22d3c3df05dd0c984042b09d5ffdd37623b.manifest.json:application/vnd.airgap.docker.distribution.manifest.v2+json \
  layers/images/postgresql/16.1.0-debian-11-r15/a0e6fe8bb3459d4fe89f14866fc85016a6fe402c645490e05b07b06ca0ed335b:application/vnd.docker.container.image.v1+json \
  layers/images/postgresql/16.1.0-debian-11-r15/8e248d9790e37ac00cdc74d5e92886f87438902fe7f09c72f6f95bda693311af:application/vnd.docker.image.rootfs.diff.tar.gzip \
  layers/images/postgresql/16.1.0-debian-11-r15/8650c8a0b51cc2b519e1c3836c6150a747878fa5fb86b2504eca9fef4213071e.manifest.json:application/vnd.airgap.docker.distribution.manifest.v2+json \
  layers/images/postgresql/16.1.0-debian-11-r15/db4d87205de2a28510f7e610676e6545b083fb8ea8671892feaf39682a739721:application/vnd.docker.container.image.v1+json \
  layers/images/postgresql/16.1.0-debian-11-r15/5e22b7f245aa934f5239b562012ca392cd58391c41676d2098e6d83d981a9ffd:application/vnd.docker.image.rootfs.diff.tar.gzip
