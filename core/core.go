package core

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Header struct {
	Key   string
	Value string
}

type Domain struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	Name         string
	ParentDomain bson.ObjectId `bson:",omitempty"`
	Skipped      bool          `bson:",omitempty"`
	LastChecked  time.Time     `bson:",omitempty"`
	Headers      []Header      `bson:",omitempty"`
}
