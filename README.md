# Airgap OCI Bundle

This repo contains some resources from initial investigations into how airgap bundles can be delivered via OCI artifacts.

## Dependencies

[oras](https://github.com/oras-project/oras) is an OCI client for managing content like artifacts, images, packages. It can be used to create an OCI artifact from a directory of files.

## Creating an OCI artifact

To login to a registry:

```bash
oras login registry.replicated.com -u $USERNAME -p $PASSWORD
```

Given the following file structure:

```                                             
layers
├── bye.txt
└── hi.txt
```

The following command will create an OCI artifact that has a custom artifact type of `application/vnd.airgap.bundle` and two layers, one with the media type `application/vnd.airgap.hi` and one with the media type `application/vnd.airgap.bye`.

```bash
oras push \
  --artifact-type application/vnd.airgap.bundle \
  registry.replicated.com/${APP_SLUG}/airgap-artifact:v1 \
  layers/hi.txt:application/vnd.airgap.hi \
  layers/bye.txt:application/vnd.airgap.bye
```

To fetch the manifest for the uploaded artifact, run:

```bash
oras manifest fetch registry.replicated.com/${APP_SLUG}/airgap-artifact:v1
```

## Checksums

An OCI image has a manifest that contains a list of layers. Each layer has a digest that is a checksum of the layer's contents. The manifest also has a digest that is a checksum of the manifest's contents. The manifest digest is used to reference the image.

To compute digest of the manifest, run:

```bash
oras manifest fetch registry.replicated.com/${APP_SLUG}/airgap-artifact:v1 | shasum -a 256
```

To pull an artifact using the manifest digest, run:

```bash
oras pull registry.replicated.com/${APP_SLUG}/airgap-artifact@sha256:${SHA256SUM}
```

If desired, the checksums of the files returned can be validated against the checksums of the layers in the manifest:

```bash
shasum -a 256 layers/*
```

## Building an Airgap Bundle

TODO - CLI example that:
* Pulls down kubernetes config (replicated release or helm chart)
* Pulls down Images (skopeo copy)
* Creates airgap metadata file (can be manual/hardcoded for now)
* Creates OCI artifact (oras push)
* Downloads the bundle (oras pull)
* Pushes images to registry (skopeo copy)
* Deploys kubernetes config
