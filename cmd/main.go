package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	xmlPath = "/home/leadness/Загрузки/ruwiki-20200701-pages-articles-multistream.xml"
)


type location struct {
	Data string `xml:",chardata"`
}

func main() {
	f, err := os.Open(xmlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	i := 0
	d := xml.NewDecoder(f)
	count := 0
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			log.Fatalf("Error decoding token: %s", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "title" {
				// If this is a start element named "location", parse this element
				// fully.
				var loc location
				if err = d.DecodeElement(&loc, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				i += 1
				fmt.Println(i, loc.Data)
				if i == 100 {
					break
				}
			}
		default:
		}
	}

}

