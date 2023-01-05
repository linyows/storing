package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

type CLI struct {
	out, err io.Writer
	// source path files
	args []string
	// logrotate is generate source path for rotation
	logrotate bool
	// bucket is name for storage
	bucket string
	// prefix is prefix of object
	prefix string
	// upload timeout
	timeout int
	// credentials data
	credentials []byte
	// store interface
	store Store
}

func (c *CLI) Do() {
	if ver {
		fmt.Fprintf(c.err, c.buildVersion()+"\n")
		return
	}

	if len(c.args) == 0 || c.bucket == "" {
		fmt.Fprintf(c.err, "Error:\n  filepath and bucket is required\n\n")
		flag.Usage()
		return
	}

	if c.credentials == nil {
		u, _ := user.Current()
		var err error
		c.credentials, err = os.ReadFile(strings.Replace(credentialsFile, "~", u.HomeDir, 1))
		if err != nil {
			fmt.Fprintf(c.err, "credentials file: %v\n", err)
			return
		}
	}

	sourcePaths, err := c.buildSourcePaths()
	if err != nil {
		fmt.Fprintf(c.err, "buildSourcePaths: %v\n", err)
		return
	}
	if c.logrotate && len(sourcePaths) == 0 {
		fmt.Fprintf(c.err, "logrotate file not found")
		return
	}

	if c.store == nil {
		c.store = &Storing{}
	}

	c.store.SetBucket(c.bucket).SetTimeout(c.timeout).SetCredentials(c.credentials)

	for _, sourcePath := range sourcePaths {
		object, err := c.buildObjectPath(sourcePath)
		if err != nil {
			fmt.Fprintf(c.err, "buildObjectPath: %v\n", err)
			return
		}

		if err := c.store.Upload(object, sourcePath); err != nil {
			fmt.Fprintf(c.err, "storing.Upload: %v\n", err)
			return
		}
		fmt.Fprintf(c.out, "Blob %v uploaded.\n", sourcePath)
	}
}

func (c *CLI) buildSourcePaths() ([]string, error) {
	var sourcePaths []string

	for _, path := range c.args {
		files, err := filepath.Glob(path)
		if err != nil {
			return sourcePaths, err
		}
		if c.logrotate {
			for _, file := range files {
				patterns := c.getLogrotatePatterns(file)
				for _, p := range patterns {
					_, err := os.Stat(p)
					if err == nil {
						sourcePaths = append(sourcePaths, p)
					}
				}
			}
		} else {
			sourcePaths = append(sourcePaths, files...)
		}
	}

	return sourcePaths, nil
}

func (c *CLI) getLogrotatePatterns(p string) []string {
	today := time.Now().Local().Format("20060102")
	return []string{
		fmt.Sprintf("%s.1", p),
		fmt.Sprintf("%s.1.gz", p),
		fmt.Sprintf("%s-%s", p, today),
		fmt.Sprintf("%s-%s.gz", p, today),
		fmt.Sprintf("%s.%s", p, today),
		fmt.Sprintf("%s.%s.gz", p, today),
	}
}

func (c *CLI) buildVersion() string {
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

func (c *CLI) buildObjectPath(localfile string) (string, error) {
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

	if c.prefix != "" {
		return fmt.Sprintf("%s%s", c.prefix, basename), nil
	}

	return fmt.Sprintf("%s/%s/%s", hostname, lastdir, basename), nil
}
