# Command line interface for Paperless-ngx

[![Latest release](https://img.shields.io/github/v/release/hansmi/papercli)][releases]
[![CI workflow](https://github.com/hansmi/papercli/actions/workflows/ci.yaml/badge.svg)](https://github.com/hansmi/papercli/actions/workflows/ci.yaml)
[![Go reference](https://pkg.go.dev/badge/github.com/hansmi/papercli.svg)](https://pkg.go.dev/github.com/hansmi/papercli)

Papercli is command line interface tool for interacting with
[Paperless-ngx][paperless], a document management system transforming physical
documents into a searchable online archive.


## Installation

[Pre-built binaries][releases]:

* Binary archives (`.tar.gz`)
* Debian/Ubuntu (`.deb`)
* RHEL/Fedora (`.rpm`)
* Microsoft Windows (`.zip`)

Docker images via GitHub's container registry:

```shell
docker pull ghcr.io/hansmi/papercli
```

With the source being available it's also possible to produce custom builds
directly using [Go][golang] or [GoReleaser][goreleaser].


[golang]: https://golang.org/
[goreleaser]: https://goreleaser.com/
[paperless]: https://docs.paperless-ngx.com/
[releases]: https://github.com/hansmi/papercli/releases/latest

<!-- vim: set sw=2 sts=2 et : -->
