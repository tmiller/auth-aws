package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Document struct {
	XMLName xml.Name    `xml:"Response"`
	Attrs   []Attribute `xml:"Assertion>AttributeStatement>Attribute"`
}

type Attribute struct {
	Name   string   `xml:"Name,attr"`
	Values []string `xml:"AttributeValue"`
}

func main() {
	xmlFile, err := ioutil.ReadFile("sample.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var q Document
	xml.Unmarshal(xmlFile, &q)

	for _, attr := range q.Attrs {
		for _, val := range attr.Values {
			fmt.Println(val)
		}
	}
}
