package main

import "encoding/xml"

type SamlResponse struct {
	XMLName xml.Name    `xml:"Response"`
	Attrs   []Attribute `xml:"Assertion>AttributeStatement>Attribute"`
}

type Attribute struct {
	Name   string   `xml:"Name,attr"`
	Values []string `xml:"AttributeValue"`
}

func parseSaml(decodedSamlResponse []byte) (samlResponse SamlResponse, err error) {
	err = xml.Unmarshal(decodedSamlResponse, &samlResponse)
	return
}
