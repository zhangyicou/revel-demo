package models

import (
	"encoding/json"
	// "github.com/revel/revel"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Params struct {
	m map[string]interface{}
}

func (p Params) Permit(keys ...string) Params {
	for k, _ := range p.m {
		found := false
		for _, kk := range keys {
			if kk == k {
				found = true
				break
			}
		}

		if !found {
			delete(p.m, k)
		}
	}

	return p
}

func MakeParams(m map[string]interface{}) Params {
	return Params{m}
}

type ArtistData struct {
	Id   bson.ObjectId `bson:"_id"  json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (a ArtistData) IsValid() bool {
	return true
}

type Artists struct {
	*mgo.Collection
}

func (c Artists) FindById(id bson.ObjectId) *Artist {
	var artist ArtistData
	err := c.FindId(id).One(&artist)
	if err != nil {
		panic(err)
	}
	return &Artist{&artist, &c}
}

func (c Artists) All() []ArtistData {
	var artists []ArtistData
	err := c.Find(nil).All(&artists)
	if err != nil {
		panic(err)
	}
	return artists
}

func (c Artists) Create(artist *ArtistData) {
	err := c.Insert(artist)
	if err != nil {
		panic(err)
	}
}

type Artist struct {
	*ArtistData
	*Artists
}

func (a Artist) save() {
	err := a.UpdateId(a.Id, a.ArtistData)
	if err != nil {
		panic(err)
	}
}

func (a Artist) Update(p Params) bool {
	bs, err := json.Marshal(p.m)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bs, a.ArtistData)
	if err != nil {
		panic(err)
	}

	if a.IsValid() { // return additionally validation result
		a.save() // TODO save what is only changed ?
		return true
	}

	return false
}
