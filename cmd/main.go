package main

import (
	"flag"
	"log"
	"os"

	"../pkg/mongodb"
	"../pkg/parser"
)

var (
	xmlPath string
	connStr string
)

func init() {
	flag.StringVar(&xmlPath, "xml-path", "", "path to wikipedia xml dump file")
	flag.StringVar(&connStr, "conn-str", "mongodb://127.0.0.1:27017", "connection string for mongodb")
}

func main() {
	flag.Parse()
	if xmlPath == "" {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(xmlPath)
	if err != nil {
		log.Printf("error on opening wikipedia xml dump file - %s\n", err)
		flag.PrintDefaults()
		return
	}
	defer file.Close()

	storage, err := mongodb.NewStorage(connStr)
	if err != nil {
		log.Printf("error on connecting to mongodb - %s\n", err)
		flag.PrintDefaults()
		return
	}

	wikiParser := parser.NewWikiParser(file, storage)
	log.Printf("start parsing %s...\n", xmlPath)
	wikiParser.Parse()
}
