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

type WikiParser struct {
	file          *os.File
	storage       *mongo.Storage
	textProcessor *WikiTextProcessor
}

func NewWikiParser(file *os.File, storage *mongo.Storage) *WikiParser {
	return &WikiParser{
		file:          file,
		storage:       storage,
		textProcessor: NewWikiTextProcessor(),
	}
}

func (p *WikiParser) Parse() {
	d := xml.NewDecoder(p.file)
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
					Text: make(map[string]mongo.ArticlePart),
				}
			} else if ty.Name.Local == "text" {
				var field xmlField
				if err = d.DecodeElement(&field, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				articleText := field.Data
				titlesSlice := p.textProcessor.GetTitles(articleText)
				for i, text := range p.textProcessor.SplitText(articleText) {
					processedText, refsSlice := p.textProcessor.ProcessText(text)
					article.Text[titlesSlice[i]] = mongo.ArticlePart{
						Text: processedText,
						Refs: refsSlice,
					}
				}
				if err := p.storage.InsertArticle(article); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					articlesCount += 1
				}
				if articlesCount % 100000 == 0 {
					log.Printf("loaded %d articles to mongo", articlesCount)
				}
			}
		default:
		}
	}
}
