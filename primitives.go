package avrogo

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type primitive struct {
	id     string
	reader func(r io.Reader) (o interface{}, err os.Error)
}

func (p primitive) Read(r io.Reader) (interface{}, os.Error) {
	return p.reader(r)
}

func readNull(r io.Reader) (interface{}, os.Error) {
	return nil, nil
}

func readBoolean(r io.Reader) (o interface{}, err os.Error) {
	p := make([]byte, 1)
	if _, readerr := io.ReadFull(r, p); readerr != nil {
		err = readerr
	} else if p[0] == 1 {
		o = true
	} else if p[0] == 0 {
		o = false
	} else {
		err = os.NewError(fmt.Sprintf("Invalid bool value: 0x%x", p[0]))
	}
	return
}

func readInt(r io.Reader) (o interface{}, err os.Error) {
	p := make([]byte, 1)
	x := uint32(0)
	for shift := uint(0); ; shift += 7 {
		if shift >= 32 {
			return nil, os.NewError("int too long!")
		} else if _, readerr := io.ReadFull(r, p); readerr != nil {
			return nil, readerr
		}
		b := uint32(p[0])
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
	}
	return int32((x >> 1) ^ uint32((int32(x&1)<<31)>>31)), nil
}

func readLong(r io.Reader) (int64, os.Error) {
	p := make([]byte, 1)
	x := uint64(0)
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 {
			return 0, os.NewError("long too long!")
		} else if _, readerr := io.ReadFull(r, p); readerr != nil {
			return 0, readerr
		}
		b := uint64(p[0])
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
	}
	return int64((x >> 1) ^ uint64((int64(x&1)<<63)>>63)), nil
}

func readILong(r io.Reader) (interface{}, os.Error) {
	return readLong(r)
}

func readFloat(r io.Reader) (interface{}, os.Error) {
	var f float32
	err := binary.Read(r, binary.LittleEndian, &f)
	return f, err
}

func readDouble(r io.Reader) (interface{}, os.Error) {
	var d float64
	err := binary.Read(r, binary.LittleEndian, &d)
	return d, err
}

func readBytes(r io.Reader) (interface{}, os.Error) {
	l, err := readLong(r)
	if err != nil {
		return nil, err
	}
	b := make([]byte, int32(l))
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readString(r io.Reader) (string, os.Error) {
	buf, err := readBytes(r)
	if err != nil {
		return "", err
	}
	return string(buf.([]byte)), nil
}

func readIString(r io.Reader) (interface{}, os.Error) {
	return readString(r)
}

var (
	primitives = map[string]Type{}
	Null       = primitive{"null", readNull}
	Boolean    = primitive{"boolean", readBoolean}
	Int        = primitive{"int", readInt}
	Long       = primitive{"long", readILong}
	Float      = primitive{"float", readFloat}
	Double     = primitive{"double", readDouble}
	Bytes      = primitive{"bytes", readBytes}
	String     = primitive{"string", readIString}
)

func init() {
	for _, t := range []primitive{Null, Boolean, Int, Long, Float, Double, Bytes, String} {
		primitives[t.id] = t
	}
}
