package avrogo

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	reader, err := os.Open("test_schema.json", os.O_RDONLY, 0444)
	if err != nil {
		t.Fatalf("Can't open file; err=%s", err)
	}
	m := Load(reader)
	if m["type"] != "record" {
		t.Fatalf("Read json incorrectly? parsed=%v", m)
	}
}
