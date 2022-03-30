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
	InsertDoc func(string, interface{}) chan statics.Result
	FindDoc   func(string, interface{}, interface{}, *options.FindOneOptions) chan error
	FindDocs  func(string, interface{}, *options.FindOptions) chan statics.Result
	UpdateDoc func(string, interface{}, interface{}, *options.UpdateOptions) chan statics.Result
}

var Api mongoApi = mongoApi{
	Constants: constants,
	InsertDoc: insertDoc,
	FindDoc:   findDoc,
	FindDocs:  findDocs,
	UpdateDoc: updateDoc,
}
