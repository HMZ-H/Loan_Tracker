package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Action    string             `bson:"action"`
	Timestamp time.Time          `bson:"timestamp"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Details   string             `bson:"details"`
}
