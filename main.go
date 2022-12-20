package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	defaultBucket          = ""
	defaultCredentialsFile = "~/.config/storing.json"
	credentialsFile        = ""
	defaultTimeout         = 50

	// flag
	prefix    string
	bucket    string
	timeout   int
	logrotate bool
	ver       bool
)

func init() {
	flag.StringVar(&prefix, "prefix", "", "object prefix (default \"<hostname>/<lastdir>/<basename>\")")
	flag.StringVar(&bucket, "bucket", defaultBucket, "bucket name")
	flag.StringVar(&credentialsFile, "key", defaultCredentialsFile, "credentials filepath")
	flag.IntVar(&timeout, "timeout", defaultTimeout, "timeout seconds for upload")
	flag.BoolVar(&logrotate, "logrotate", false, "search today's file for logrotate")
	flag.BoolVar(&ver, "version", false, "show build version")

	helpText := `Usage:
  storing [options] <filepath...>

The options are:
`
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), helpText)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	cli := &CLI{
		out:       os.Stdout,
		err:       os.Stderr,
		args:      flag.Args(),
		bucket:    bucket,
		prefix:    prefix,
		logrotate: logrotate,
		store:     &Storing{timeout: time.Duration(timeout) * time.Second},
	}

	cli.Do()
}
