package importer

import (
	"strings"
	"testing"
)

func TestParseRows(t *testing.T) {
	rows, warnings, err := ParseRows(strings.NewReader("r1,A,2,2\nbad\nr2,B,x,3\n"))
	if err != nil || len(rows) != 1 || len(warnings) != 2 {
		t.Fatalf("rows %#v warnings %#v err %v", rows, warnings, err)
	}
}
