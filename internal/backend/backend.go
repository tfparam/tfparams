// Package backend fetches a plan JSON from a URI, dispatching on the scheme
// (s3://, gs://, azblob://, or a local path). Cloud backends are stubbed until
// Track C-3; local paths are fully supported.
package backend

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Source is a parsed --env target.
type Source struct {
	Scheme    string // "s3", "gs", "azblob", or "" for a local path
	Bucket    string // s3/gs
	Container string // azblob
	Account   string // azblob (optional, may come from AZURE_STORAGE_ACCOUNT)
	Key       string // object key
	Path      string // local path
	Raw       string
}

// Parse splits a URI into its components. A value with no scheme is a local path.
func Parse(uri string) (Source, error) {
	s := Source{Raw: uri}
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" {
		s.Path = uri
		return s, nil
	}
	s.Scheme = u.Scheme
	key := strings.TrimPrefix(u.Path, "/")
	switch u.Scheme {
	case "s3", "gs":
		s.Bucket, s.Key = u.Host, key
		if s.Bucket == "" || s.Key == "" {
			return s, fmt.Errorf("invalid %s URI %q (want %s://bucket/key)", u.Scheme, uri, u.Scheme)
		}
	case "azblob":
		s.Container, s.Key = u.Host, key
		if u.User != nil {
			s.Account = u.User.Username()
		}
		if s.Container == "" || s.Key == "" {
			return s, fmt.Errorf("invalid azblob URI %q (want azblob://[account@]container/key)", uri)
		}
	default:
		return s, fmt.Errorf("unsupported URI scheme %q", u.Scheme)
	}
	return s, nil
}

// Fetch reads the bytes the URI points at, dispatching on scheme. Cloud
// credentials come from each SDK's default credential chain.
func Fetch(ctx context.Context, uri string) ([]byte, error) {
	s, err := Parse(uri)
	if err != nil {
		return nil, err
	}
	switch s.Scheme {
	case "":
		return fetchLocal(s)
	case "s3":
		return fetchS3(ctx, s)
	case "gs":
		return fetchGCS(ctx, s)
	case "azblob":
		return fetchAzblob(ctx, s)
	default:
		return nil, fmt.Errorf("unsupported URI scheme %q", s.Scheme)
	}
}
