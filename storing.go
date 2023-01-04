package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Store interface {
	SetBucket(b string) Store
	SetTimeout(t int) Store
	SetCredentials(c []byte) Store
	Upload(dst, src string) error
}

type Storing struct {
	bucket      string
	timeout     time.Duration
	credentials []byte
}

func (s *Storing) SetBucket(b string) Store {
	s.bucket = b
	return s
}

func (s *Storing) SetTimeout(t int) Store {
	s.timeout = time.Duration(t) * time.Second
	return s
}

func (s *Storing) SetCredentials(c []byte) Store {
	s.credentials = c
	return s
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
