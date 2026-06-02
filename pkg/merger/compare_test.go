package merger

import "testing"

func findRow(rows []CompareRow, name string) (CompareRow, bool) {
	for _, r := range rows {
		if r.Name == name {
			return r, true
		}
	}
	return CompareRow{}, false
}

func TestCompare(t *testing.T) {
	dev := EnvResult{Name: "dev", Params: []Param{
		{Name: "instance_type", Description: "type", HasApplied: true, Applied: "t3.small"},
		{Name: "region", HasApplied: true, Applied: "ap-northeast-1"},
		{Name: "db_password", Sensitive: true, HasApplied: true, Applied: "dev-pw"},
	}}
	prd := EnvResult{Name: "prd", Params: []Param{
		{Name: "instance_type", Description: "type", HasApplied: true, Applied: "t3.xlarge"},
		{Name: "region", HasApplied: true, Applied: "ap-northeast-1"},
		{Name: "db_password", Sensitive: true, HasApplied: true, Applied: "prd-pw"},
	}}

	rows := Compare([]EnvResult{dev, prd})

	it, _ := findRow(rows, "instance_type")
	if !it.Differs {
		t.Errorf("instance_type should differ: %+v", it)
	}
	if it.Values["dev"] != "t3.small" || it.Values["prd"] != "t3.xlarge" {
		t.Errorf("instance_type values = %+v", it.Values)
	}

	region, _ := findRow(rows, "region")
	if region.Differs {
		t.Errorf("region should NOT differ: %+v", region)
	}

	pw, _ := findRow(rows, "db_password")
	if !pw.Sensitive || pw.Differs {
		t.Errorf("db_password should be sensitive and not diff-highlighted: %+v", pw)
	}
}

func TestCompareMissingEnvIsNotSet(t *testing.T) {
	dev := EnvResult{Name: "dev", Params: []Param{{Name: "x", HasApplied: true, Applied: "1"}}}
	prd := EnvResult{Name: "prd", Params: []Param{}} // x missing in prd
	rows := Compare([]EnvResult{dev, prd})
	x, ok := findRow(rows, "x")
	if !ok {
		t.Fatal("x missing")
	}
	if x.Values["prd"] != "(not set)" {
		t.Errorf("prd x should be (not set), got %q", x.Values["prd"])
	}
	if !x.Differs {
		t.Errorf("x should differ (1 vs not set)")
	}
}
