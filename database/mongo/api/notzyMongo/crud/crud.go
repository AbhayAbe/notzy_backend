package crud

import (
	"context"

	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertDoc(collectionName string, doc interface{}) chan utils.Result {
	var db *mongo.Database = database.DB
	ch := make(chan utils.Result)
	go func() {
		defer close(ch)
		_, err := db.Collection(collectionName).InsertOne(context.Background(), doc)
		if err != nil {
			ch <- utils.Result{Error: err, Result: nil}
		} else {
			ch <- utils.Result{Error: nil, Result: doc}
		}
	}()
	return ch
}
func FindDoc(collectionName string, filter interface{}, decode interface{}, opts *options.FindOneOptions) chan error {
	var db *mongo.Database = database.DB
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

func UpdateDoc(collectionName string, filter interface{}, update interface{}, opts *options.UpdateOptions) chan utils.Result {
	var db *mongo.Database = database.DB
	ch := make(chan utils.Result)
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
			ch <- utils.Result{Error: err, Result: nil}
		} else {
			ch <- utils.Result{Error: nil, Result: res}
		}
	}()
	return ch
}
