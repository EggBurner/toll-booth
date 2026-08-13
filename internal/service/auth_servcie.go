package service

import (
	"context"
	"log"
	"net/mail"
	"strings"
	"toll-booth/internal/db"
	"unicode"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func IsValidEmail(email string) (string, bool) {
	_, err := mail.ParseAddress(email)

	if err != nil {
		log.Printf("Error pasring mail to check if its valid %v", err)
		return "Invalid email " + err.Error(), false
	}

	collection := db.Client.Database("toll-booth-db").Collection("users")

	filter := bson.M{
		"email": email,
	}

	err = collection.FindOne(context.TODO(), filter).Err()

	if err == mongo.ErrNoDocuments {
		return "successful", true
	} else if err != nil {
		return "database error: " + err.Error(), false
	} else {
		return "email already exists", false
	}

}

func IsValidWIthoutDBCheck(email string) (string, bool) {
	_, err := mail.ParseAddress(email)

	if err != nil {
		log.Printf("Error pasring mail to check if its valid %v", err)
		return "Invalid email " + err.Error(), false
	}

	return "Valid email", true
}

func IsValidPassword(password string) (string, bool) {
	if len(password) < 8 {
		return "password too short", false
	}
	if len(password) > 70 {
		return "pasword too long", false
	}
	var hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(`!@#$%^&*(),.?":{}|<>`, r):
			hasSpecial = true
		}
	}
	return "incorrect password - lacking a special character/digit/uppercase", hasUpper && hasDigit && hasSpecial
}
