package avrogo

import (
	"fmt"
	"io"
	"json"
	"os"
)

type SchemaError struct {
	msg string
}

func (e SchemaError) String() string {
	return e.msg
}

type Type interface {
	Read(r io.Reader) (o interface{}, err os.Error)
}

type Field struct {
	name     string
	doc      string
	ftype    Type
	defvalue interface{}
}

type Record struct {
	name      string
	namespace string
	doc       string
	aliases   []string
	fields    []Field
}

func (r Record) Id() string {
	if r.namespace == "" {
		return r.name
	}
	return r.namespace + "." + r.name
}

func (rec Record) Read(r io.Reader) (o interface{}, err os.Error) {
	vals := make(map[string]interface{})
	for _, f := range rec.fields {
		if val, err := f.ftype.Read(r); err == nil {
			vals[f.name] = val
		} else {
			return nil, err
		}
	}
	return vals, nil
}

type Map struct {
	values Type
}

func (m Map) Read(r io.Reader) (interface{}, os.Error) {
	vals := make(map[string]interface{})
	for count, err := readLong(r); count != 0; count, err = readLong(r) {
		if err != nil {
			return nil, err
		}
		if count < 0 {
			if _, err := readLong(r); err != nil { // Ignore the size
				return nil, err
			}
			count = -count
		}
		for ; count > 0; count-- {
			if k, err := readString(r); err != nil {
				return nil, err
			} else if v, err := m.values.Read(r); err != nil {
				return nil, err
			} else {
				vals[k] = v
			}
		}
	}
	return vals, nil
}

func getString(obj map[string]interface{}, name string) string {
	if v, ok := obj[name]; !ok {
		return ""
	} else if s, ok := v.(string); ok {
		return s
	}
	panic("errrrrr")
}

func requireString(obj map[string]interface{}, name string) string {
	if v, ok := obj[name]; !ok {
		panic(SchemaError{"Missing required field '" + name + "'"})
	} else if s, ok := v.(string); ok {
		return s
	}
	panic(SchemaError{name + " must be a string"})

}

func getStringArray(obj map[string]interface{}, name string) []string {
	if v, ok := obj[name]; !ok {
		return []string{}
	} else if a, ok := v.([]string); ok {
		return a
	}
	panic(SchemaError{name + " must be an array of string, not " + obj[name].(string)})
}

func loadField(obj map[string]interface{}) Field {
	return Field{requireString(obj, "name"), getString(obj, "doc"), loadType(obj["type"]),
		getString(obj, "default")}
}

func loadMap(obj map[string]interface{}) Map {
	return Map{loadType(requireString(obj, "values"))}
}

func loadRecord(obj map[string]interface{}) Record {
	var fields []Field
	for _, v := range obj["fields"].([]interface{}) {
		fields = append(fields, loadField(v.(map[string]interface{})))
	}
	return Record{requireString(obj, "name"), getString(obj, "namespace"), getString(obj, "doc"),
		getStringArray(obj, "aliases"), fields}
}

func loadType(i interface{}) Type {
	switch v := i.(type) {
	case string:
		if p, ok := primitives[v]; ok {
			return p
		} else {
			panic(SchemaError{"Unknown type name " + v})
		}
	case []interface{}:
		panic(SchemaError{"Not handling unions yet!"})
	case map[string]interface{}:
		if t, ok := v["type"]; !ok {
			panic(SchemaError{fmt.Sprintf("Type object without type name %s", v)})
		} else if t == "map" {
			return loadMap(v)
		} else if t == "record" {
			return loadRecord(v)
		}
	}
	panic(SchemaError{"Unknown type: " + i.(string)})
}

func Load(r io.Reader) (loaded Type, err os.Error) {
	var i interface{}
	d := json.NewDecoder(r)
	d.Decode(&i)
	defer func() {
		if e := recover(); e != nil {
			loaded = nil
			err = e.(SchemaError)
		}
	}()
	loaded = loadType(i)
	return
}
