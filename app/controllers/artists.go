package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"revel-demo/app/models"
)

func cacheModel(f func() interface{}, keys ...interface{}) interface{} {
	var keyStrs []string
	for _, key := range keys {
		keyStrs = append(keyStrs, fmt.Sprint(key))
	}

	key := strings.Join(keyStrs, "-")

	var value interface{}
	if err := cache.Get(key, value); err != nil {
		value = f()
		go cache.Set(key, value, 1*time.Minute)
	}

	return value
}

func assert(e error) {
	if e != nil {
		panic(e)
	}
}

type Artists struct {
	App
}

func (c Artists) Recover() {
	switch err := recover(); err {
	case mgo.ErrNotFound:
		c.Result = c.NotFound("Could not be found")
	case ErrNotAuthorized:
		c.Result = c.Forbidden("Not allowed")
	}
}

var ErrNotAuthorized = errors.New("not authorized")

func (c Artists) Authorize() {
	switch c.MethodName {
	case "Create":
		// some logic based on current user and other elements
		if false {
			panic(ErrNotAuthorized)
		}
	}
}

func (c Artists) collection() models.Artists {
	return models.Artists{c.mongoSession.DB("martys").C("artists")}

}

func (c Artists) Show(id bson.ObjectId) revel.Result {
	defer c.Recover()

	artist := cacheModel(func() interface{} {
		return c.collection().FindById(id)
	}, "artist", id).(*models.Artist)

	return c.RenderJson(artist.ArtistData)
}

func (c Artists) Index() revel.Result {
	defer c.Recover()

	artists := c.collection().All()

	return c.RenderJson(artists)
}

func (c Artists) ReadModel(model interface{}) {
	body, err := ioutil.ReadAll(c.Request.Body)
	assert(err)

	err = json.Unmarshal(body, model)
	assert(err)
}

func (c Artists) Create() revel.Result {
	defer c.Recover()

	c.Authorize()

	var artist *models.ArtistData
	c.ReadModel(artist)

	if artist.IsValid() {
		c.collection().Create(artist)
		return c.RenderJson(artist)
	} else {
		// validation failed
		return nil
	}
}

func (c Artists) Params() models.Params {
	var m map[string]interface{}

	body, err := ioutil.ReadAll(c.Request.Body)
	assert(err)

	err = json.Unmarshal(body, m)
	assert(err)

	return models.MakeParams(m)
}

func (c Artists) ArtistParams() models.Params {
	return c.Params().Permit(
		"name",
	)
}

func (c Artists) Update(id bson.ObjectId) revel.Result {
	defer c.Recover()

	// c.Authorize()

	artist := c.collection().FindById(id)

	// c.Authorize(artist)

	if artist.Update(c.ArtistParams()) {
		// success
		return nil
	} else {
		// validation failed
		return nil
	}
}
