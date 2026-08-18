package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"toll-booth/internal/auth"
	"toll-booth/internal/db"
	"toll-booth/internal/model"
	"toll-booth/internal/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Approved bool   `json:"approved"`
}

type ResetPasswordReuest struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Approved    bool   `json:"approved"`
}

func HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	reason, validEmail := service.IsValidEmail(req.Email)

	if !validEmail {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}

	reasonPass, validPassword := service.IsValidPassword(req.Password)

	if !validPassword {
		http.Error(w, reasonPass, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		http.Error(w, "server error, try again!", http.StatusInternalServerError)
		return
	}

	user := model.User{
		UserID:           bson.NewObjectID(),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		DateSignUp:       time.Now(),
		ActiveLinksCount: 0,
		TotalLinksCount:  0,
		Password:         hashedPassword,
		Role:             model.RoleEnum(req.Role),
	}

	collection := db.Client.
		Database("toll-booth-db").
		Collection("users")

	result, err := collection.InsertOne(
		context.TODO(),
		user,
	)

	if err != nil {
		log.Printf("failed to insert user: %v", err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"id":      result.InsertedID,
	})
}

func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, isValidEmail := service.IsValidWIthoutDBCheck(req.Email)

	_, isValidPassword := service.IsValidPassword(req.Password)

	if !(isValidEmail && isValidPassword) {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	var user model.User

	collection := db.Client.Database("toll-booth-db").Collection("users")

	filter := bson.M{
		"email": req.Email,
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "User not exists", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		err = auth.ComparePassword(user.Password, req.Password)

		if err != nil {
			http.Error(w, "Invalid Credentials", http.StatusBadRequest)
			return
		} else {
			req.Approved = true
		}
	}

	if req.Approved {
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":        "Credentials valid user logged in",
			"approvalStatus": req.Approved,
		})
	}
}

func HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ResetPasswordReuest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	_, isValidEmail := service.IsValidWIthoutDBCheck(req.Email)
	_, isValidOldPassword := service.IsValidPassword(req.OldPassword)
	_, isValidNewPassword := service.IsValidPassword(req.NewPassword)

	if !(isValidEmail && isValidOldPassword) {
		http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return
	}

	if !isValidNewPassword {
		http.Error(w, "Invalid New Password", http.StatusBadRequest)
		return
	}

	var user model.User

	collection := db.Client.Database("toll-booth-db").Collection("users")

	filter := bson.M{
		"email": req.Email,
	}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "No Record Found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error while reading the database", http.StatusInternalServerError)
		return
	} else {
		err = auth.ComparePassword(user.Password, req.OldPassword)
		if err != nil {
			http.Error(w, "Invalid Credentials", http.StatusBadRequest)
			return
		} else {
			req.Approved = true
		}
	}

	newPasswordHashed, err := auth.HashPassword(req.NewPassword)

	if err != nil {
		http.Error(w, "Server Error, try again", http.StatusInternalServerError)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"password": newPasswordHashed,
		},
	}

	if req.Approved {
		result, err := collection.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			http.Error(w, "Error updating db", http.StatusInternalServerError)
		}

		log.Printf("Update found %v", result.MatchedCount)
		log.Printf("Update modified %v", result.ModifiedCount)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Password Updated",
		})
	}

}
