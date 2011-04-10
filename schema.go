package avrogo

import (
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
	Id() string
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
		return loadRecord(v)
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
