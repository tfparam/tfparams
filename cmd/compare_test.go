package cmd

import (
	"strings"
	"testing"
)

func TestCompareCommand(t *testing.T) {
	out, err := runCmd(t, "compare",
		"--env", "dev=../testdata/plan_dev.json",
		"--env", "prd=../testdata/plan_prd.json",
		"--docs-json", "../testdata/docs.json",
	)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{
		"# Environment Comparison",
		"| Name | Description | dev | prd | Diff |",
		"`t3.small`",
		"`t3.xlarge`",
		"⚠️",          // instance_type / replica_count differ
		"(sensitive)", // db_password masked
	} {
		if !strings.Contains(out, want) {
			t.Errorf("compare output missing %q\n%s", want, out)
		}
	}
}

func TestCompareRequiresTwoEnvs(t *testing.T) {
	_, err := runCmd(t, "compare", "--env", "dev=../testdata/plan_dev.json", "--docs-json", "../testdata/docs.json")
	if err == nil {
		t.Fatal("expected error: compare needs >= 2 envs")
	}
}

func TestCompareCloudBackendNotImplemented(t *testing.T) {
	_, err := runCmd(t, "compare",
		"--env", "dev=../testdata/plan_dev.json",
		"--env", "prd=s3://bucket/prd/plan.json",
		"--docs-json", "../testdata/docs.json",
	)
	if err == nil {
		t.Fatal("expected not-implemented error for s3 backend")
	}
}
