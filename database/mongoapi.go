package database

import (
	"github.com/AbhayAbe/notzy_backend/statics"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConstants struct {
	IsUnique string
}

var constants mongoConstants = mongoConstants{
	IsUnique: "isUnique",
}

type mongoApi struct {
	Constants mongoConstants
	InsertDoc func(collection string, doc interface{}) chan statics.Result
	FindDoc   func(collection string, filter interface{}, decode interface{}, options *options.FindOneOptions) chan error
	FindDocs  func(collection string, filter interface{}, options *options.FindOptions) chan statics.Result
	UpdateDoc func(collection string, filter interface{}, update interface{}, options *options.UpdateOptions) chan statics.Result
	DeleteDoc func(collection string, filter interface{}, options *options.DeleteOptions) chan statics.Result
}

var Api mongoApi = mongoApi{
	Constants: constants,
	InsertDoc: insertDoc,
	FindDoc:   findDoc,
	FindDocs:  findDocs,
	UpdateDoc: updateDoc,
	DeleteDoc: deleteDoc,
}
