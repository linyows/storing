package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	version       = "dev"
	commit        = ""
	date          = ""
	builtBy       = ""
	defaultBucket = ""
	credsPath     = "creds/gcp.json"

	// flag
	object string
	bucket string
	ver    bool

	//go:embed creds
	creds embed.FS
)

type Storing struct {
	bucket  string
	timeout time.Duration
	creds   []byte
}

func init() {
	flag.StringVar(&object, "object", "", "object name (default format \"hostname/lastdir/basename\")")
	flag.StringVar(&bucket, "bucket", defaultBucket, "bucket name")
	flag.BoolVar(&ver, "version", false, "show build version")

	helpText := `Usage:
  storing [options] <filepath>

The options are:
`
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), helpText)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if ver {
		fmt.Fprintf(os.Stderr, buildVersion(version, commit, date, builtBy)+"\n")
		return
	}
	localfile := os.Args[1]

	if localfile == "" || bucket == "" {
		if localfile == "" {
			fmt.Fprintf(os.Stderr, "localfile was required\n")
		}
		if bucket == "" {
			fmt.Fprintf(os.Stderr, "bucket name was required\n")
		}
		return
	}

	if object == "" {
		op, err := buildObjectPath(localfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "buildObjectPath: %v\n", err)
			return
		}
		object = op
	}

	json, err := creds.ReadFile(credsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "embed.FS Readfile: %v\n", err)
		return
	}

	s := &Storing{
		bucket:  bucket,
		timeout: time.Second * 50,
		creds:   json,
	}
	if err := s.Upload(object, localfile); err != nil {
		fmt.Fprintf(os.Stderr, "storing.Upload: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "Blob %v uploaded.\n", localfile)
}

func buildObjectPath(localfile string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("os.Hostname: %v", err)
	}

	abspath, err := filepath.Abs(localfile)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: %v", err)
	}

	basename := filepath.Base(abspath)
	dirname := filepath.Dir(abspath)
	dirs := strings.Split(dirname, "/")
	lastdir := dirs[len(dirs)-1]

	return fmt.Sprintf("%s/%s/%s", hostname, lastdir, basename), nil
}

func (s *Storing) Upload(dst, src string) error {
	// Create clinet
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(s.creds))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Create writer and Upload
	o := client.Bucket(s.bucket).Object(dst)
	o = o.If(storage.Conditions{DoesNotExist: true})
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}

	return result
}
