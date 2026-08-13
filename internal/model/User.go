package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoleEnum string

const (
	UserMember RoleEnum = "MEMBER"
	UserAdmin  RoleEnum = "ADMIN"
	UserOwner  RoleEnum = "OWNER"
)

type User struct {
	UserID           bson.ObjectID `bson:"_id"`
	FirstName        string        `bson:"firstName"`
	LastName         string        `bson:"lastName"`
	Email            string        `bson:"email"`
	Password         string        `bson:"password"`
	DateSignUp       time.Time     `bson:"dateSignUp"`
	ActiveLinksCount uint16        `bson:"activeLinksCount"`
	TotalLinksCount  uint32        `bson:"totalLinksCount"`
	Role             RoleEnum      `bson:"role"`
}
