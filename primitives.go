package avrogo

import (
	"fmt"
	"io"
	"os"
)

type primitive struct {
	id     string
	reader func(r io.Reader) (o interface{}, err os.Error)
}

func (p primitive) Id() string {
	return p.id
}

func (p primitive) Read(r io.Reader) (interface{}, os.Error) {
	return p.reader(r)
}


func readNull(r io.Reader) (interface{}, os.Error) {
	return nil, nil
}

var Null = primitive{"null", readNull}

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

var Boolean = primitive{"boolean", readBoolean}

func readInt(r io.Reader) (o interface{}, err os.Error) {
	p := make([]byte, 1)
	x := uint32(0)
	for shift := uint(0); ; shift += 7 {
		// TODO - bail if we go over the int length
		if _, readerr := io.ReadFull(r, p); readerr != nil {
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

var Int = primitive{"int", readInt}

type Long struct{}

func (l Long) Id() string {
	return "long"
}

type Float struct{}

func (f Float) Id() string {
	return "float"
}

type Double struct{}

func (d Double) Id() string {
	return "double"
}

type Bytes struct{}

func (b Bytes) Id() string {
	return "bytes"
}

type String struct{}

func (s String) Id() string {
	return "string"
}

var primitives = map[string]Type{}

func init() {
	for _, t := range []Type{Null, Boolean, Int} { //, Long{}, Float{}, Double{}, Bytes{}, String{}} {
		primitives[t.Id()] = t
	}
}
