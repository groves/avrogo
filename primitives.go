package avrogo

import (
	"io"
	"os"
)

type Null struct{}

func (n Null) Id() string {
	return "null"
}

func (n Null) Read(r io.Reader) (o interface{}, err os.Error) {
	return nil, nil
}

type Boolean struct{}

func (b Boolean) Id() string {
	return "boolean"
}
func (b Boolean) Read(r io.Reader) (o interface{}, err os.Error) {
	p := make([]byte, 1)
	if _, err := io.ReadFull(r, p); err != nil {
		return nil, err
	}
	return p[0] == 1, nil
}

type Int struct{}

func (i Int) Id() string {
	return "int"
}

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
	for _, t := range []Type{Null{}, Boolean{}} { //Int{}, Long{}, Float{}, Double{}, Bytes{}, String{}} {
		primitives[t.Id()] = t
	}
}
