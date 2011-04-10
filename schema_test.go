package avrogo

import (
	"bytes"
	"os"
	"testing"
)

func load(t *testing.T, fn string) Type {
	loaded, err := loadErr(t, fn)
	if err != nil {
		t.Fatalf("Schema error: %s", err)
	}
	return loaded
}

func loadErr(t *testing.T, fn string) (Type, os.Error) {
	reader, err := os.Open(fn, os.O_RDONLY, 0444)
	if err != nil {
		t.Fatalf("Unable to open file %s", fn)
	}
	return Load(reader)
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

func TestNamelessRecord(t *testing.T) {
	_, err := loadErr(t, "nameless_record_schema.json")
	if err.String() != "Missing required field 'name'" {
		t.Fatalf("Incorrect SchemaError: %s", err)
	}
}
