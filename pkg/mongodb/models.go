package mongodb

import (
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	ID    bson.ObjectId          `bson:"_id" json:"_id,omitempty"`
	Title string                 `bson:"title" json:"title"`
	Text  map[string]ArticlePart `bson:"text" json:"text"`
}

type ArticlePart struct {
	Text string   `bson:"text" json:"text"`
	Refs []string `bson:"refs" json:"refs"`
}
