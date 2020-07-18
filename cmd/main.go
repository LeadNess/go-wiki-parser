package main

import (
	"log"

	"../pkg/parser"
)

const (
	xmlPath = "/home/leadness/Загрузки/ruwiki-20200701-pages-articles-multistream.xml"
	connStr = "mongodb://172.17.0.2:27017"
)

func main() {
	if err := parser.ParseWikiXml(xmlPath, connStr); err != nil {
		log.Fatal(err)
	}
}
