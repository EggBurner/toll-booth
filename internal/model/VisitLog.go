package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VisitLog struct {
	VisitLogID  bson.ObjectID `bson:"_id"`
	LinkID      bson.ObjectID `bson:"linkID"`
	TimeStamp   time.Time     `bson:"timestamp"`
	IPAddress   string        `bson:"ipAddress"`
	Location    string        `bson:"location"`
	Referrer    string        `bson:"referrer"`
	MachineType string        `bson:"machineType"`
}
