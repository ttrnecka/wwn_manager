package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

const GLOBAL_CUSTOMER = "__GLOBAL__"
const UNKNOWN_CUSTOMER = "<NO CUSTOMER>"

func NilObjectID() primitive.ObjectID {
	return primitive.NilObjectID
}

func NilOjectIdSlice() []primitive.ObjectID {
	return make([]primitive.ObjectID, 0)
}
