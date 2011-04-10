package avrogo

import (
	"os"
	"testing"
)

func load(t *testing.T, fn string) Type {
	reader, err := os.Open(fn, os.O_RDONLY, 0444)
	if err == nil {
		return Load(reader)
	}
	t.Fatalf("Unable to open file %s", fn)
	panic("Shouldn't get here")
}

func SkipTestDecode(t *testing.T) {
	m := load(t, "test_schema.json")
	if m.Id() != "org.apache.avro.Interop" {
		t.Fatalf("Read json incorrectly? parsed=%v", m)
	}
}

func TestPrimitiveDecode(t *testing.T) {
	m := load(t, "primitive_record_schema.json")
	if m.Id() != "test.AllPrimitives" {
		t.Fatalf("Read record id incorrectly?%v", m)
	}
	if len(m.(Record).fields) != 2 {
		t.Fatalf("Got incorrect number of fields?%v", m)
	}

}
