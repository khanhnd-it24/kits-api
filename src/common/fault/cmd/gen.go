package main

import (
	"bytes"
	"go/format"
	"gopkg.in/yaml.v2"
	strutil "kits/api/src/common/string"
	"log"
	"os"
	"text/template"
)

type ErrDetail struct {
	Code    int64
	Message string
}

type ErrMessage struct {
	CamelKey string
	Key      string
	Val      ErrDetail
}

type TemplateData struct {
	Errors []ErrMessage
}

func main() {
	yFile, err := os.ReadFile("key.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]ErrDetail)
	err = yaml.Unmarshal(yFile, &data)

	if err != nil {
		log.Fatal(err)
	}

	errorKeys := make([]ErrMessage, 0, len(data))
	for k, v := range data {
		errorKeys = append(errorKeys, ErrMessage{strutil.ToCamel(k), k, v})
	}

	temp, err := template.ParseFiles("key.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	err = temp.Execute(&buf, TemplateData{errorKeys})
	if err != nil {
		log.Fatal(err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fileGen, err := os.Create("keygen.go")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fileGen.Write(p)
	if err != nil {
		log.Fatal(err)
	}
}
