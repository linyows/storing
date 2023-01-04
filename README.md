Storing
==

The storing is the cloud storage upload CLI.
It may support AWS S3 if needed, but for now only GCS is supported.

  <a href="https://github.com/linyows/storing/actions/workflows/test.yml"><img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/linyows/storing/test.yml?branch=main&label=Test&style=for-the-badge"></a>
  <a href="https://github.com/linyows/storing/actions/workflows/build.yml"><img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/linyows/storing/build.yml?branch=main&style=for-the-badge"></a>
  <a href="https://github.com/linyows/storing/releases"><img src="http://img.shields.io/github/release/linyows/storing.svg?style=for-the-badge" alt="GitHub Release"></a>
  <a href="https://github.com/linyows/storing/blob/main/LICENSE"><img src="http://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge" alt="MIT License"></a>
  <a href="http://godoc.org/github.com/linyows/storing"><img src="http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge" alt="Go Documentation"></a>
  <a href="https://codecov.io/gh/linyows/storing"> <img src="https://img.shields.io/codecov/c/github/linyows/storing.svg?style=for-the-badge" alt="codecov"></a>

Usage
--

```sh
$ go install github.com/linyows/storing@latest
$ storing ./testdata/example.jpg -bucket my-bucket-name -key ~/Downloads/my-project-credentials.yml
Blob ./testdata/example.jpg uploaded: https://storage.cloud.google.com/<my-bucket-name>/<hostname>/testdata/example.jpg
```

Options
--

The options are:

```
  -bucket string
        bucket name
  -key string
        credentials filepath (default "~/.config/storing.json")
  -logrotate
        search today's file for logrotate
  -prefix string
        object prefix (default "<hostname>/<lastdir>/<basename>")
  -timeout int
        timeout seconds for upload (default 50)
```

Author
--

[linyows](https://github.com/linyows)
