Storing
==

The storing is the cloud storage upload CLI which contains the credentials for upload.
It may support AWS S3 if needed, but for now only GCS is supported.

Usage
--

In order to use the credential embedded in the go binary, place the credential in the `creds` directory, pass the bucket name as a variable and go build.

```sh
$ git clone https://github.com/linyows/storing.git && cd storing
$ mv ~/Downloads/my-project-credentials.yml ./creds/gcp.json
$ BUCKET=my-bucket-name make
$ ./storing ./testdata/example.jpg
Blob ./testdata/example.jpg uploaded: https://storage.cloud.google.com/<my-bucket-name>/<hostname>/testdata/example.jpg
```

Options
--

The options are:

```
  -bucket string
        bucket name (default "hosting-service-proxylog")
  -object string
        object name (default format "hostname/lastdir/basename")
```

Author
--

[linyows](https://github.com/linyows)
