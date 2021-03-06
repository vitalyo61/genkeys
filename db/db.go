package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/vitalyo61/genkeys/db/model"
)

type DB struct {
	session        *mgo.Session
	name           string
	codeCollection string
}

func Make(url string) (*DB, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &DB{
		session:        session,
		name:           "gencode",
		codeCollection: "codes",
	}, nil
}

func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) CodeSet(code *model.Code) error {
	conn := db.session.Copy()
	coll := conn.DB(db.name).C(db.codeCollection)
	defer conn.Close()

	return coll.Insert(code)
}

func (db *DB) CodeGet(number string) (*model.Code, error) {
	conn := db.session.Copy()
	coll := conn.DB(db.name).C(db.codeCollection)
	defer conn.Close()

	c := new(model.Code)
	err := coll.FindId(number).One(c)
	return c, err
}

func (db *DB) CodeExtinguish(number string) error {
	conn := db.session.Copy()
	coll := conn.DB(db.name).C(db.codeCollection)
	defer conn.Close()

	return coll.Update(
		bson.M{
			"_id":    bson.M{"$eq": number},
			"status": bson.M{"$eq": model.CodeYes},
		},
		bson.M{
			"$set": bson.M{"status": model.CodeStop},
		})
}

func (db *DB) CodeLast() (*model.Code, error) {
	conn := db.session.Copy()
	coll := conn.DB(db.name).C(db.codeCollection)
	defer conn.Close()

	c := new(model.Code)
	err := coll.Find(nil).Sort("-_id").One(c)
	return c, err
}

func (db *DB) CodeRemove(number string) error {
	conn := db.session.Copy()
	coll := conn.DB(db.name).C(db.codeCollection)
	defer conn.Close()

	return coll.RemoveId(number)
}
