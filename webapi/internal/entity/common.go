package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

const GlobalCustomer = "__GLOBAL__"
const UnknownCustomer = "<NO CUSTOMER>"

type ID = primitive.ObjectID

func NilObjectID() primitive.ObjectID {
	return primitive.NilObjectID
}

func NilOjectIDSlice() []primitive.ObjectID {
	return make([]primitive.ObjectID, 0)
}
