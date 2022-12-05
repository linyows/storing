Storing
==

The storing is the cloud storage upload CLI.
It may support AWS S3 if needed, but for now only GCS is supported.

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
  -object string
        object name (default format "<hostname>/<lastdir>/<basename>")
```

Author
--

[linyows](https://github.com/linyows)
