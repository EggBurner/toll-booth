package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"toll-booth/internal/auth"
	"toll-booth/internal/db"
	"toll-booth/internal/model"
	"toll-booth/internal/service"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LinkShortenRequest struct {
	TargetLink   string `json:"targetLink"`
	OwnerID      string `json:"ownerID"`
	PinProtected bool   `json:"pinProtected"`
	PIN          string `json:"PIN"`
}

func HandleLinkShortner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LinkShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body "+err.Error(), http.StatusBadRequest)
		return
	}

	isValidURL := service.IsValidURL(req.TargetLink)
	var hashedPin string

	if !(isValidURL) {
		http.Error(w, "Invalid Request, data format uncorrect", http.StatusBadRequest)
		return
	}
	if req.PinProtected {
		isValidPin := service.IsValidPIN(req.PIN)
		if !isValidPin {
			http.Error(w, "Invalid Request, data format uncorrect", http.StatusBadRequest)
			return
		}
		hashedPin, err = auth.HashPassword(req.PIN)
		if err != nil {
			http.Error(w, "Internal server error, try again "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	ownerID, err := bson.ObjectIDFromHex(req.OwnerID)
	if err != nil {
		http.Error(w, "Invalid owner ID", http.StatusBadRequest)
		return
	}

	shortCode, err := service.LinkShortner()

	if err != nil {
		http.Error(w, "Error while generating short code", http.StatusInternalServerError)
		return
	}

	link := model.Link{
		LinkID:          bson.NewObjectID(),
		TargetLink:      req.TargetLink,
		ShortCode:       shortCode,
		LinkDateCreated: time.Now(),
		OwnerID:         ownerID,
		Status:          model.StatusEnum("ACTIVE"),
		VisitCount:      0,
		PinProtected:    req.PinProtected,
	}

	if req.PinProtected {
		link.PinHash = &hashedPin
	}

	collection := db.Client.Database("toll-booth-db").Collection("links")

	result, err := collection.InsertOne(context.TODO(), link)

	if err != nil {
		http.Error(w, "Could not register link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Link created successfully",
		"id":      result.InsertedID,
	})

}
