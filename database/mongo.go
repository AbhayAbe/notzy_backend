package database

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/AbhayAbe/notzy_backend/database/mongo/constants"
	"github.com/AbhayAbe/notzy_backend/models"
	"github.com/AbhayAbe/notzy_backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _URI_ string
var Client *mongo.Client
var DB *mongo.Database

func ConfigureMongodb() {

	var err error
	_URI_ = os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(_URI_)
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	if err := Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	DB = Client.Database("notzy")

	res := <-initIndices(models.User{})
	fmt.Println(res.Result)
}

func DisconnectMongodb() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func initIndices(model interface{}) <-chan utils.Result {
	ch := make(chan utils.Result)
	go func() {
		defer close(ch)
		mod := reflect.TypeOf(model)
		schemaName := strings.ToLower(mod.Name()) + "s"
		for i := 0; i < mod.NumField(); i++ {
			field := mod.Field(i)
			tag := field.Tag.Get(constants.IsUnique)
			key := field.Tag.Get("json")
			if len(tag) > 0 && len(key) > 0 {
				_, err := DB.Collection(schemaName).Indexes().CreateOne(context.Background(),
					mongo.IndexModel{
						Keys:    bson.M{key: 1},
						Options: options.Index().SetUnique(true),
					})
				if err != nil {
					DisconnectMongodb()
					fmt.Println("Error: ", err)
					ch <- utils.Result{Error: err, Result: 0}
				}
				fmt.Println("Index for", schemaName, "created")
				ch <- utils.Result{Error: nil, Result: 1}
			}
		}
	}()
	return ch
}
