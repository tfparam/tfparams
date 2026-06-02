package parser

import (
	"strings"
	"testing"
)

func TestParseDocs(t *testing.T) {
	in := `{"inputs":[
      {"name":"instance_type","type":"string","description":"EC2","default":"t3.medium","required":false},
      {"name":"db_password","type":"string","description":"pw","default":null,"required":true,"sensitive":true}
    ]}`
	d, err := ParseDocs(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseDocs: %v", err)
	}
	if len(d.Inputs) != 2 {
		t.Fatalf("inputs = %d, want 2", len(d.Inputs))
	}
	if !d.Inputs[0].HasDefault() || FormatValue(d.Inputs[0].Default) != "t3.medium" {
		t.Errorf("instance_type default mismatch: %s", string(d.Inputs[0].Default))
	}
	if d.Inputs[1].HasDefault() {
		t.Errorf("db_password should have no default, got %s", string(d.Inputs[1].Default))
	}
	if !d.Inputs[1].Required || !d.Inputs[1].Sensitive {
		t.Errorf("db_password should be required and sensitive")
	}
}
