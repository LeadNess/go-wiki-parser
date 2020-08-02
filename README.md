# go-wiki-parser

[![Build Status](https://travis-ci.com/LeadNess/go-wiki-parser.svg?branch=master)](https://travis-ci.com/LeadNess/go-wiki-parser)

### Description

Parses large xml wiki dump into MongoDB.

### Usage

- ```git clone https://github.com/LeadNess/go-wiki-parser.git```
- ```go get -t github.com/LeadNess/go-wiki-parser```
- ```go get -t go-wiki-parser/cmd/main.go```
- ```go run go-wiki-parser/cmd/main.go -xml-path <path-to-xml-dump-file> -conn-str <connection-string-for-mongo>```