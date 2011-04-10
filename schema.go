package avrogo

import (
	"io"
	"json"
	"os"
)

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
		if val, err := f.ftype.Read(r) ; err == nil {
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

func getStringArray(obj map[string]interface{}, name string) []string {
	if v, ok := obj[name]; !ok {
		return []string{}
	} else if a, ok := v.([]string); ok {
		return a
	}
	panic("errrr")
}

func loadField(obj map[string]interface{}) Field {
	return Field{obj["name"].(string), getString(obj, "doc"), loadType(obj["type"]),
		getString(obj, "default")}
}

func loadRecord(obj map[string]interface{}) Record {
	var fields []Field
	for _, v := range obj["fields"].([]interface{}) {
		fields = append(fields, loadField(v.(map[string]interface{})))
	}
	return Record{obj["name"].(string), getString(obj, "namespace"), getString(obj, "doc"),
		getStringArray(obj, "aliases"), fields}
}

func loadType(i interface{}) Type {
	switch v := i.(type) {
	case string:
		if p, ok := primitives[v]; ok {
			return p
		} else {
			panic("Unknown type name " + v)
		}
	case []interface{}:
		panic("Not handling unions yet!")
	case map[string]interface{}:
		return loadRecord(v)
	}
	panic("Unknown type: " + i.(string))
}

func Load(r io.Reader) Type {
	var i interface{}
	d := json.NewDecoder(r)
	d.Decode(&i)
	return loadType(i)
}
