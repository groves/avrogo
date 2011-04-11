package avrogo

import (
	"bytes"
	"testing"
)

func read(t *testing.T, ty Type, data []byte, expected interface{}) {
	if v, err := ty.Read(bytes.NewBuffer(data)); err != nil {
		t.Errorf("Error reading %x: %v", data, err)
	} else if v != expected {
		t.Errorf("Expected '%v' for '%x' but got '%v'", expected, data, ty, v)
	}
}

func readBadValue(t *testing.T, ty Type, data []byte) {
	if v, err := ty.Read(bytes.NewBuffer(data)); err == nil {
		t.Errorf("Expected error reading '%x', but read '%v' successfully", v, data)
	}
}

func TestBooleanEncoding(t *testing.T) {
	read(t, Boolean, []byte{0x01}, true)
	read(t, Boolean, []byte{0x00}, false)
	readBadValue(t, Boolean, []byte{0x02})
}

func TestIntEncoding(t *testing.T) {
	read(t, Int, []byte{0x00}, int32(0))
	read(t, Int, []byte{0x01}, int32(-1))
	read(t, Int, []byte{0x02}, int32(1))
	read(t, Int, []byte{0x7f}, int32(-64))
	read(t, Int, []byte{0x80, 0x01}, int32(64))
	read(t, Int, []byte{0x81, 0x01}, int32(-65))
}

func TestLongEncoding(t *testing.T) {
	read(t, Long, []byte{0x00}, int64(0))
	read(t, Long, []byte{0x01}, int64(-1))
	read(t, Long, []byte{0x80, 0x01}, int64(64))
	read(t, Long, []byte{0x81, 0x01}, int64(-65))
}

func TestFloatEncoding(t *testing.T) {
	read(t, Float, []byte{0x00, 0x00, 0x00, 0x00}, float32(0.0))
}

func TestDoubleEncoding(t *testing.T) {
	read(t, Double, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, float64(0.0))
}
