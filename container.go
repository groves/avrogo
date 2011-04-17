package avrogo

import (
	"bytes"
	"io"
	"os"
)

type Iterator struct {
	r                io.Reader
	schema           Type
	remainingInBlock int64
	next             interface{}
	err              os.Error
}

func (i *Iterator) Next() (ret interface{}) {
	if i.next == nil {
		if hasNext, err := i.HasNext(); err != nil {
			panic(err)
		} else if !hasNext {
			panic("Iterator.Next called without a call to HasNext")
		}
	}
	ret = i.next
	i.next = nil
	return
}

func (i *Iterator) HasNext() (bool, os.Error) {
	for i.err == nil && i.next == nil {
		if i.remainingInBlock == 0 {
			if i.remainingInBlock, i.err = readLong(i.r); i.err != nil {
				return false, i.err
			} else if _, i.err = readLong(i.r); i.err != nil { // Skip block object byte size
				return false, i.err
			}
		}
		if i.next, i.err = i.schema.Read(i.r); i.err != nil {
			return false, i.err
		}
		i.remainingInBlock--
		if i.remainingInBlock == 0 {
			// Read the sync marker
			_, i.err = io.ReadFull(i.r, make([]byte, 16)) // TODO - check against the sync marker
		}
	}
	return i.next != nil, i.err
}


var headerSchema = map[string]interface{}{
	"type": "record", "name": "org.apache.avro.file.Header", "fields": []interface{}{
		map[string]interface{}{"name": "magic", "type": map[string]interface{}{"type": "fixed", "name": "Magic", "size": 4}},
		map[string]interface{}{"name": "meta", "type": map[string]interface{}{"type": "map", "values": "bytes"}},
		map[string]interface{}{"name": "sync", "type": map[string]interface{}{"type": "fixed", "name": "Sync", "size": 16}}},
}
var headerRecord Type

func init() {
	var err os.Error
	if headerRecord, err = LoadSchema(headerSchema); err != nil {
		panic(err)
	}
}

func ReadContainer(r io.Reader) (schema Type, iter *Iterator, err os.Error) {
	if header, headererr := headerRecord.Read(r); headererr != nil {
		err = headererr
	} else if metabytes, metabytesok := header.(map[string]interface{})["meta"]; !metabytesok {
		err = SchemaError{"Header didn't contain a meta record"}
	} else if schemabytes, schemabytesok := metabytes.(map[string]interface{})["avro.schema"]; !schemabytesok {
		err = SchemaError{"Metadata didn't contain 'avro.schema'"}
	} else if schema, err = LoadJsonSchema(bytes.NewBuffer(schemabytes.([]byte))); err == nil {
		iter = &Iterator{r: r, schema: schema}
	}
	return

}
