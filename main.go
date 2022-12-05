package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	version                = "dev"
	commit                 = ""
	date                   = ""
	builtBy                = ""
	defaultBucket          = ""
	defaultCredentialsFile = "~/.config/storing.json"
	credentialsFile        = ""

	// flag
	object string
	bucket string
	ver    bool
)

func init() {
	flag.StringVar(&object, "object", "", "object name (default format \"<hostname>/<lastdir>/<basename>\")")
	flag.StringVar(&bucket, "bucket", defaultBucket, "bucket name")
	flag.StringVar(&credentialsFile, "key", defaultCredentialsFile, "credentials filepath")
	flag.BoolVar(&ver, "version", false, "show build version")

	helpText := `Usage:
  storing [options] <filepath>

The options are:
`
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), helpText)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	if ver {
		fmt.Fprintf(os.Stderr, buildVersion(version, commit, date, builtBy)+"\n")
		return
	}

	if len(args) == 0 {
		flag.Usage()
		return
	}

	localfile := args[0]
	if localfile == "" || bucket == "" {
		fmt.Fprintf(os.Stderr, "Error:\n  filepath and bucket is required\n\n")
		flag.Usage()
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

	u, _ := user.Current()
	json, err := os.ReadFile(strings.Replace(credentialsFile, "~", u.HomeDir, 1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "credentials file: %v\n", err)
		return
	}

	s := &Storing{
		bucket:      bucket,
		timeout:     time.Second * 50,
		credentials: json,
	}
	if err := s.Upload(object, localfile); err != nil {
		fmt.Fprintf(os.Stderr, "storing.Upload: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "Blob %v uploaded.\n", localfile)
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

type Storing struct {
	bucket      string
	timeout     time.Duration
	credentials []byte
}

func (s *Storing) Upload(dst, src string) error {
	// Create clinet
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(s.credentials))
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
