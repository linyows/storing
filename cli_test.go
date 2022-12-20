package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetLogrotatePatterns(t *testing.T) {
	cli := &CLI{}
	today := time.Now().Local().Format("20060102")
	want := []string{
		"access_log.1",
		"access_log.1.gz",
		fmt.Sprintf("access_log-%s", today),
		fmt.Sprintf("access_log-%s.gz", today),
		fmt.Sprintf("access_log.%s", today),
		fmt.Sprintf("access_log.%s.gz", today),
	}
	got := cli.getLogrotatePatterns("access_log")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %s, but got %s", want, got)
	}
}

func TestBuildSourcePaths(t *testing.T) {
	cli := &CLI{
		out:  os.Stdout,
		err:  os.Stderr,
		args: []string{"testdata/*.jpg"},
	}
	want := []string{
		"testdata/example.jpg",
	}
	got, err := cli.buildSourcePaths()
	if err != nil {
		t.Errorf("buildSourcePaths error: %#v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %s, but got %s", want, got)
	}
}

func TestBuildObjectPath(t *testing.T) {
	hostname, _ := os.Hostname()
	localfile := "testdata/example.jpg"
	cli := &CLI{}
	want := fmt.Sprintf("%s/testdata/example.jpg", hostname)
	got, err := cli.buildObjectPath(localfile)
	if err != nil {
		t.Errorf("buildObjectPath error: %#v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %s, but got %s", want, got)
	}
}

type MStoring struct {
	bucket      string
	timeout     time.Duration
	credentials []byte
	dst         string
	src         string
}

func (s *MStoring) SetBucket(b string) Store {
	s.bucket = b
	return s
}

func (s *MStoring) SetTimeout(t time.Duration) Store {
	s.timeout = t
	return s
}

func (s *MStoring) SetCredentials(c []byte) Store {
	s.credentials = c
	return s
}

func (s *MStoring) Upload(dst, src string) error {
	s.dst = dst
	s.src = src
	return nil
}

func TestDo(t *testing.T) {
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	bucket := "linyows"
	storing := &MStoring{
		timeout: 10 * time.Second,
	}
	cli := &CLI{
		out:       stdout,
		err:       stderr,
		args:      []string{"testdata/*.jpg"},
		logrotate: false,
		bucket:    bucket,
		store:     storing,
	}
	cli.Do()

	hostname, _ := os.Hostname()

	wantDst := fmt.Sprintf("%s/testdata/example.jpg", hostname)
	if storing.dst != wantDst {
		t.Errorf("want %s, but got %s", wantDst, storing.dst)
	}
	wantSrc := "testdata/example.jpg"
	if storing.src != wantSrc {
		t.Errorf("want %s, but got %s", wantSrc, storing.src)
	}
}
