package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

const GLOBAL_CUSTOMER = "__GLOBAL__"

func NilObjectID() primitive.ObjectID {
	return primitive.NilObjectID
}
