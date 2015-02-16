package controllers

import (
	"fmt"
	"reflect"

	"github.com/revel/revel"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"revel-demo/app/models"
)

var (
	mongoSession *mgo.Session
)

type App struct {
	mongoSession *mgo.Session
	*revel.Controller
	User *models.User
}

func (c App) Index() revel.Result {
	return c.Render()
}

// TODO where do it ?
func init() {
	objId := bson.NewObjectId()
	revel.TypeBinders[reflect.TypeOf(objId)] = ObjectIdBinder

	var err error
	mongoSession, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	revel.InterceptMethod((*App).SetUser, revel.BEFORE)
	revel.InterceptMethod((*App).MongoBefore, revel.BEFORE)
	revel.InterceptMethod((*App).MongoFinally, revel.FINALLY)
}

var ObjectIdBinder = revel.Binder{
	// Make a ObjectId from a request containing it in string format.
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		if len(val) == 0 {
			return reflect.Zero(typ)
		}
		if bson.IsObjectIdHex(val) {
			objId := bson.ObjectIdHex(val)
			return reflect.ValueOf(objId)
		} else {
			revel.ERROR.Print("ObjectIdBinder.Bind - invalid ObjectId!")
			return reflect.Zero(typ)
		}
	}),
	// Turns ObjectId back to hexString for reverse routing
	Unbind: func(output map[string]string, name string, val interface{}) {
		var hexStr string
		hexStr = fmt.Sprintf("%s", val.(bson.ObjectId).Hex())
		// not sure if this is too carefull but i wouldn't want invalid ObjectIds in my App
		if bson.IsObjectIdHex(hexStr) {
			output[name] = hexStr
		} else {
			revel.ERROR.Print("ObjectIdBinder.Unbind - invalid ObjectId!")
			output[name] = ""
		}
	},
}

func (c *App) MongoBefore() revel.Result {
	c.mongoSession = mongoSession.Clone()
	return nil
}

func (c App) MongoFinally() revel.Result {
	c.mongoSession.Close()
	return nil
}

func (c *App) SetUser() revel.Result {
	// TODO find user by token
	c.User = &models.User{}
	return nil
}

func (c App) Recover() revel.Result {
	if e := recover(); e != nil {
		fmt.Println(e)
	}
	return nil
}
