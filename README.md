Storing
==

The storing is the cloud storage upload CLI which contains the credentials for upload.
If you don't embed a credential, storing will read `~/.config/storing.json`.
It may support AWS S3 if needed, but for now only GCS is supported.

⚠️ Notice: When embedding a credential, be careful with go's binary handling. Also, regardless of embedding the credential, the credential should have the least privileges (creator privileges for one bucket).

Usage
--

It describes two types of use.

Standalone:

```sh
$ go install github.com/linyows/storing@latest
$ storing ./testdata/example.jpg -bucket my-bucket-name -key ~/Downloads/my-project-credentials.yml
Blob ./testdata/example.jpg uploaded: https://storage.cloud.google.com/<my-bucket-name>/<hostname>/testdata/example.jpg
```

Bucket and credentials embed:

In order to use the credential embedded in the go binary, place the credential in the `embed` directory, pass the bucket name as a variable and go build.

```sh
$ git clone https://github.com/linyows/storing.git && cd storing
$ mv ~/Downloads/my-project-credentials.yml ./embed/credentials.json
$ BUCKET=my-bucket-name make
$ ./storing ./testdata/example.jpg
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
