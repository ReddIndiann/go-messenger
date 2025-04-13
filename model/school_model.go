package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	Verified  bool               `bson:"verified" json:"verified"`
	Logo      string             `bson:"logo" json:"logo"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
