package avrogo

import (
	"bytes"
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
	r := m.(Record)
	if len(r.fields) != 2 {
		t.Fatalf("Got incorrect number of fields?%v", m)
	}
	encoded := []byte{0x01}
	decoded, _ := r.Read(bytes.NewBuffer(encoded))
	if decoded.(map[string]interface{})["nullField"] != nil {
		t.Fatalf("Got wrong type for nullField")
	}
	if decoded.(map[string]interface{})["boolField"] != true {
		t.Fatalf("Got wrong value for boolField")
	}


}
