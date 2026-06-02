package backend

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		uri     string
		scheme  string
		bucket  string
		cont    string
		account string
		key     string
		path    string
		wantErr bool
	}{
		{uri: "s3://my-bucket/env/prd/plan.json", scheme: "s3", bucket: "my-bucket", key: "env/prd/plan.json"},
		{uri: "gs://b/k/plan.json", scheme: "gs", bucket: "b", key: "k/plan.json"},
		{uri: "azblob://acct@cont/prd/plan.json", scheme: "azblob", cont: "cont", account: "acct", key: "prd/plan.json"},
		{uri: "azblob://cont/plan.json", scheme: "azblob", cont: "cont", key: "plan.json"},
		{uri: "./plan.json", scheme: "", path: "./plan.json"},
		{uri: "/abs/plan.json", scheme: "", path: "/abs/plan.json"},
		{uri: "s3://only-bucket", wantErr: true},
	}
	for _, c := range cases {
		s, err := Parse(c.uri)
		if c.wantErr {
			if err == nil {
				t.Errorf("Parse(%q): expected error", c.uri)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q): %v", c.uri, err)
			continue
		}
		if s.Scheme != c.scheme || s.Bucket != c.bucket || s.Container != c.cont ||
			s.Account != c.account || s.Key != c.key || s.Path != c.path {
			t.Errorf("Parse(%q) = %+v", c.uri, s)
		}
	}
}

func TestFetchLocal(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "plan.json")
	if err := os.WriteFile(p, []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	data, err := Fetch(context.Background(), p)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("got %q", data)
	}
}

func TestFetchLocalMissing(t *testing.T) {
	_, err := Fetch(context.Background(), filepath.Join(t.TempDir(), "nope.json"))
	if err == nil {
		t.Fatal("expected error for missing local file")
	}
}

func TestFetchUnsupportedScheme(t *testing.T) {
	if _, err := Fetch(context.Background(), "ftp://host/path"); err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
}
