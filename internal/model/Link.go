package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StatusEnum string

const (
	StatusActive   StatusEnum = "ACTIVE"
	StatusExpired  StatusEnum = "EXPIRED"
	StatusDisabled StatusEnum = "DISABLED"
)

type Link struct {
	LinkID          bson.ObjectID `bson:"_id"`
	TargetLink      string        `bson:"targetLink"`
	ShortCode       string        `bson:"shortCode"`
	LinkDateCreated time.Time     `bson:"linkDateCreated"`
	OwnerID         bson.ObjectID `bson:"ownerID"`
	Status          StatusEnum    `bson:"status"`
	VisitCount      uint32        `bson:"visitCount"`
	PinProtected    bool          `bson:"pinProtected"`
	PinHash         *string       `bson:"pinHash,omitempty"`
}
