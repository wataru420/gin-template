package util

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"reflect"
)

type SX interface{}

// go binary encoder
func ToGOB64(m interface{}) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string, m interface{}) SX {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		log.Println(`failed gob Decode`, err)
		log.Println(`type`, reflect.TypeOf(m))
	}
	return m
}
