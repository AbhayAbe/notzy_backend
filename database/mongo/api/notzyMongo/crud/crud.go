package crud

import (
	"context"

	"github.com/AbhayAbe/notzy_backend/database"
	"github.com/AbhayAbe/notzy_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDoc(collectionName string, user interface{}) chan utils.Result {
	var db *mongo.Database = database.DB
	ch := make(chan utils.Result)
	go func() {
		defer close(ch)
		_, err := db.Collection(collectionName).InsertOne(context.Background(), user)
		if err != nil {
			ch <- utils.Result{Error: err, Result: nil}
		} else {
			ch <- utils.Result{Error: nil, Result: user}
		}
	}()
	return ch
}
