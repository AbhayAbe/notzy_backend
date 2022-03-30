package database

import (
	"context"

	"github.com/AbhayAbe/notzy_backend/statics"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Insert
func insertDoc(collectionName string, doc interface{}) chan statics.Result {
	var db *mongo.Database = DB
	ch := make(chan statics.Result)
	go func() {
		defer close(ch)
		res, err := db.Collection(collectionName).InsertOne(context.Background(), doc)
		if err != nil {
			ch <- statics.Result{Error: err, Result: nil}
		} else {
			ch <- statics.Result{Error: nil, Result: gin.H{"_id": res.InsertedID.(primitive.ObjectID).Hex(), "data": doc}}
		}
	}()
	return ch
}

//Find
func findDoc(collectionName string, filter interface{}, decode interface{}, opts *options.FindOneOptions) chan error {
	var db *mongo.Database = DB
	ch := make(chan error)
	go func() {
		defer close(ch)
		var err error
		if opts != nil {
			err = db.Collection(collectionName).FindOne(context.Background(), filter, opts).Decode(decode)
		} else {
			err = db.Collection(collectionName).FindOne(context.Background(), filter).Decode(decode)
		}

		if err != nil {
			ch <- err
		} else {
			ch <- nil
		}
	}()
	return ch
}

func findDocs(collectionName string, filter interface{}, opts *options.FindOptions) chan statics.Result {
	var db *mongo.Database = DB
	ch := make(chan statics.Result)
	go func() {
		defer close(ch)
		var cur *mongo.Cursor
		var err error
		if opts != nil {
			cur, err = db.Collection(collectionName).Find(context.Background(), filter, opts)
		} else {
			cur, err = db.Collection(collectionName).Find(context.Background(), filter)
		}

		if err != nil {
			ch <- statics.Result{Error: err, Result: nil}
		} else {
			ch <- statics.Result{Error: nil, Result: cur}
		}
	}()
	return ch
}

//Update
func updateDoc(collectionName string, filter interface{}, update interface{}, opts *options.UpdateOptions) chan statics.Result {
	var db *mongo.Database = DB
	ch := make(chan statics.Result)
	go func() {
		defer close(ch)
		var res interface{}
		var err error
		if opts != nil {
			res, err = db.Collection(collectionName).UpdateOne(context.Background(), filter, update, opts)
		} else {
			res, err = db.Collection(collectionName).UpdateOne(context.Background(), filter, update)
		}

		if err != nil {
			ch <- statics.Result{Error: err, Result: nil}
		} else {
			ch <- statics.Result{Error: nil, Result: res}
		}
	}()
	return ch
}
