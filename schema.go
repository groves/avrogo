package avrogo

import (
	"io"
	"json"
)

type Type interface {
	Id() string
}

type Field struct {
	name     string
	doc      string
	ftype    interface{}
	defvalue interface{}
}

type Record struct {
	name      string
	namespace string
	doc       string
	aliases   []string
}

func (r Record) Id() string {
	if r.namespace == "" {
		return r.name
	}
	return r.namespace + "." + r.name
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
	return Field{obj["name"].(string), getString(obj, "doc"), obj["type"],
		getString(obj, "default")}
}

func loadRecord(obj map[string]interface{}) Record {
	var fields []Field
	for _, v := range obj["fields"].([]interface{}) {
		fields = append(fields, loadField(v.(map[string]interface{})))
	}
	return Record{obj["name"].(string), getString(obj, "namespace"), getString(obj, "doc"),
		getStringArray(obj, "aliases")}
}

func loadType(obj map[string]interface{}) Type {
	switch t, _ := obj["type"].(string); t {
	case "record":
		return loadRecord(obj)
	}
	panic("Unknown type: " + obj["type"].(string))
}

func Load(r io.Reader) Type {
	var i interface{}
	d := json.NewDecoder(r)
	d.Decode(&i)
	switch v := i.(type) {
	case string:
		// TODO lookup name
	case []interface{}:
		// TODO load union
	case map[string]interface{}:
		return loadType(v)
	}
	panic("Unhandled type")
}
