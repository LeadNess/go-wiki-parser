package parser

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	mongo "../mongodb"
)


type xmlField struct {
	Data string `xml:",chardata"`
}

func ParseWikiXml(xmlFilePath, connStr string) error {
	f, err := os.Open(xmlFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	storage, err := mongo.NewStorage(connStr)
	if err != nil {
		return err
	}

	parser := NewWikiParser()
	d := xml.NewDecoder(f)
	articlesCount := 0

	var article mongo.Article
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error decoding token: %s", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "title" {
				var field xmlField
				if err = d.DecodeElement(&field, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				article = mongo.Article{
					Title: field.Data,
				}
			} else if ty.Name.Local == "text" {
				var field xmlField
				if err = d.DecodeElement(&field, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				articleText := field.Data
				titlesSlice := parser.GetTitles(articleText)
				for i, text := range parser.SplitText(articleText) {
					processedText, refsSlice := parser.ProcessText(text)
					article.Text[titlesSlice[i]] = mongo.ArticlePart{
						Text: processedText,
						Refs: refsSlice,
					}
				}
				if err := storage.InsertArticle(article); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					articlesCount += 1
				}
				if articlesCount % 100000 == 0 {
					log.Printf("loaded to mongo %d articles", articlesCount)
				}
			}
		default:
		}
	}
}
