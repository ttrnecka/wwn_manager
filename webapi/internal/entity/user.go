package entity

import (
	cdb "github.com/ttrnecka/agent_poc/common/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	cdb.BaseModel `bson:",inline"`
	Username      string `bson:"username"`
	Email         string `bson:"email"`
	// #nosec G117
	Password string `bson:"password,omitempty"`
}

func Users(db *mongo.Database) *cdb.CRUD[User] {
	return cdb.NewCRUD[User](db, "users")
}
