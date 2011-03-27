package avrogo

import (
	"io"
	"json"
)

func Load(r io.Reader) map[string]interface{} {
	m := make(map[string]interface{})
	d := json.NewDecoder(r)
	d.Decode(&m)
	return m
}
