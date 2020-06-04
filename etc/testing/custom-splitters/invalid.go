package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

func Split(buf io.Reader, out chan []byte) {
	decoder := xml.NewDecoder(buf)

	for {
		token, err := decoder.Token()
		if err != nil {
			panic(err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			out <- []byte(fmt.Sprintf("+ %s", t.Name.Local))
		case xml.EndElement:
			out <- []byte(fmt.Sprintf("- %s", t.Name.Local))
		default:
		}
	}
}