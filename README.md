<!--
SPDX-FileCopyrightText: 2024 Shun Sakai

SPDX-License-Identifier: Apache-2.0 OR MIT
-->

# lzip-go

[![CI][ci-badge]][ci-url]
[![Go Reference][reference-badge]][reference-url]
![Go version][go-version-badge]

**lzip-go** is an implementation of the [lzip compressed format].

## Usage

To install this library:

```sh
go get -u github.com/sorairolake/lzip-go
```

### Example

Please see [`example_test.go`].

### Documentation

See the [documentation][reference-url] for more details.

## Minimum Go version

This library requires the minimum version of Go 1.21.

## Changelog

Please see [CHANGELOG.adoc].

## Contributing

Please see [CONTRIBUTING.adoc].

## License

Copyright &copy; 2024 Shun Sakai (see [AUTHORS.adoc])

This library is distributed under the terms of either the _Apache License 2.0_
or the _MIT License_.

This project is compliant with version 3.0 of the [_REUSE Specification_]. See
copyright notices of individual files for more details on copyright and
licensing information.

[ci-badge]: https://img.shields.io/github/actions/workflow/status/sorairolake/lzip-go/CI.yaml?branch=develop&style=for-the-badge&logo=github&label=CI
[ci-url]: https://github.com/sorairolake/lzip-go/actions?query=branch%3Adevelop+workflow%3ACI++
[reference-badge]: https://img.shields.io/badge/Go-Reference-steelblue?style=for-the-badge&logo=go
[reference-url]: https://pkg.go.dev/github.com/sorairolake/lzip-go
[go-version-badge]: https://img.shields.io/github/go-mod/go-version/sorairolake/lzip-go?style=for-the-badge&logo=go
[lzip compressed format]: https://www.nongnu.org/lzip/manual/lzip_manual.html#File-format
[`example_test.go`]: example_test.go
[CHANGELOG.adoc]: CHANGELOG.adoc
[CONTRIBUTING.adoc]: CONTRIBUTING.adoc
[AUTHORS.adoc]: AUTHORS.adoc
[_REUSE Specification_]: https://reuse.software/spec/
