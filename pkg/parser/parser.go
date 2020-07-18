package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)
type Article struct {
	Title string
	Text  map[string]interface{}
}


type xmlField struct {
	Data string `xml:",chardata"`
}

func ParseWikiXml(xmlPath string) error {
	f, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error decoding token: %s", err)
		}

		parser := NewWikiParser()
		var article Article
		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "title" {
				var field xmlField
				if err = d.DecodeElement(&field, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				article.Title = field.Data
			} else if ty.Name.Local == "text" {
				var field xmlField
				if err = d.DecodeElement(&field, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				articleText := field.Data
				titlesSlice := parser.GetTitles(articleText)
				for i, text := range parser.SplitText(articleText) {
					processedText, refsSlice := parser.ProcessText(text)
					article.Text[titlesSlice[i]] = map[string]interface{}{
						"text": processedText,
						"refs": refsSlice,
					}
				}
			}
		default:
		}
	}
}
